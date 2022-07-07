package persistence

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type connection struct {
	inner *driver.Conn
	auth  auth
}

func (s *server) addr() string {
	return s.host + ":" + strconv.FormatInt(int64(s.port), 10)
}

func (s *server) options(a *auth) *clickhouse.Options {
	return &clickhouse.Options{
		Addr: []string{s.addr()},
		Auth: clickhouse.Auth{
			Database: s.db,
			Username: a.user,
			Password: a.pwd,
		},
		Debug:           true,
		ConnMaxLifetime: time.Hour,
	}
}

func (s *server) open(a *auth) *driver.Conn {
	conn, err := clickhouse.Open(s.options(a))
	if err != nil {
		log.Printf("clickhouse connection failed: %s", err)
		return nil
	}
	return &conn
}

func (s *server) save(conn *connection, data string, name string) {
	query := fmt.Sprintf(
		`INSERT INTO %s SETTINGS input_format_skip_unknown_fields = true FORMAT JSONEachRow
			%s`, name, data)
	err := (*conn.inner).AsyncInsert(s.ctx, query, false)
	if err != nil {
		log.Fatalf("clickhouse insert failed: %+v", err)
	}
}

// START - Used for integration tests. Keeping ti dirty for now.
// @todo: Refactor or move big part of this to testing utility class.

func (s *server) RecycleTables(conn *connection) {

	dropQuery := `DROP TABLE IF EXISTS integration_test_webperf_rum_events`

	dropErr := (*conn.inner).Exec(s.ctx, dropQuery)

	if dropErr != nil {
		log.Fatal(dropErr)
	}

	createQuery := `CREATE TABLE IF NOT EXISTS integration_test_webperf_rum_events (
		event_date Date DEFAULT toDate(created_at),
		created_at DateTime,
		event_type                      LowCardinality(String),
		browser_name                    LowCardinality(String),
		browser_version                 String,
		device_manufacturer             LowCardinality(String),
		device_type                     LowCardinality(String),
		user_agent                      String,
		next_hop_protocol               LowCardinality(String),
		visibility_state                LowCardinality(String),
	
		session_id                      FixedString(43),
		session_length                  UInt8,
		url                             String,
		connect_duration                Nullable(UInt16),
		dns_duration                    Nullable(UInt16),
		first_byte_duration             Nullable(UInt16),
		redirect_duration               Nullable(UInt16),
		redirects_count                 UInt8,
		
		first_contentful_paint          Nullable(UInt16),
		first_paint                     Nullable(UInt16),
	
		cumulative_layout_shift         Nullable(Float32),
		first_input_delay               Nullable(UInt16),
		largest_contentful_paint        Nullable(UInt16),
	
		country_code                    FixedString(2),

		boomerang_version               LowCardinality(String),
		screen_width                    Nullable(UInt16),
		screen_height                   Nullable(UInt16),

		dom_res                         Nullable(UInt16),
		dom_doms                        Nullable(UInt16),
		mem_total                       Nullable(UInt32),
		mem_limit                       Nullable(UInt32),
		mem_used                        Nullable(UInt32),
		mem_lsln                        Nullable(UInt32),
		mem_ssln                        Nullable(UInt32),
		mem_lssz                        Nullable(UInt32)
	)
		ENGINE = MergeTree()
		PARTITION BY toYYYYMMDD(event_date)
		ORDER BY (device_type, event_date)
		SETTINGS index_granularity = 8192`

	createErr := (*conn.inner).Exec(s.ctx, createQuery)

	if createErr != nil {
		log.Fatal(createErr)
	}
}

func (s *server) countRecords(conn *connection) {
	rows, err := (*conn.inner).Query(s.ctx, "SELECT count(*) FROM integration_test_webperf_rum_events")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var (
			col1 uint64
		)
		if err := rows.Scan(&col1); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("row: Count=%d\n", col1)
	}
	rows.Close()
}

// END - Used for integration tests.
