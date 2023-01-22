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
		Debug:           false,
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

func (s *server) save(conn *connection, data string, table string) {
	if data != "" && table != "" {
		query := fmt.Sprintf(
			`INSERT INTO %s SETTINGS input_format_skip_unknown_fields = true FORMAT JSONEachRow
				%s`, table, data)
		err := (*conn.inner).AsyncInsert(s.ctx, query, false)
		if err != nil {
			log.Printf("clickhouse insert failed: %+v", err)
		}
	} else {
		log.Printf("clickhouse invalid data for table %s: %s", table, data)
	}
}

// CheckTableExist checks if table exists
func (s *server) CheckTableExist(conn *connection, table string) (bool, error) {
	if table == "" {
		return false, fmt.Errorf("table is required")
	}
	query := fmt.Sprintf(`EXISTS %s`, table)
	rows, err := (*conn.inner).Query(s.ctx, query)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		var result *uint8
		err = rows.Scan(&result)
		if err != nil {
			return false, err
		}
		return result != nil && *result == 1, nil
	}
	return false, fmt.Errorf("no rows found")
}

// CreateTable creates the table if not exists
func (s *server) CreateTable(conn *connection, table string) error {
	createQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		event_date                      Date DEFAULT toDate(created_at),
		hostname                        LowCardinality(String),
		created_at                      DateTime,
		event_type                      LowCardinality(String),
		browser_name                    LowCardinality(String),
		browser_version                 Nullable(String),
		ua_vnd                          LowCardinality(Nullable(String)),
		ua_plt                          LowCardinality(Nullable(String)),
		device_type                     LowCardinality(String),
		device_manufacturer             LowCardinality(Nullable(String)),
		operating_system                LowCardinality(String),
		operating_system_version        Nullable(String),
		user_agent                      Nullable(String),
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
	
		geo_country_code                FixedString(2),
		geo_city_name                   Nullable(String),
		page_id                         FixedString(8),

		data_saver_on                   Nullable(UInt8),

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
		mem_lssz                        Nullable(UInt32),
		scr_bpp                         Nullable(String),
		scr_orn                         Nullable(String),
		cpu_cnc                         Nullable(UInt8),
		dom_ln                          Nullable(UInt16),
		dom_sz                          Nullable(UInt16),
		dom_ck                          Nullable(UInt16),
		dom_img                         Nullable(UInt16),
		dom_img_uniq                    Nullable(UInt16),
		dom_script                      Nullable(UInt16),
		dom_iframe                      Nullable(UInt16),
		dom_link                        Nullable(UInt16),
		dom_link_css                    Nullable(UInt16)
	)
		ENGINE = MergeTree()
		PARTITION BY toYYYYMMDD(event_date)
		ORDER BY (hostname, event_date)
		SETTINGS index_granularity = 8192`, table)

	log.Printf("creating table with query: %v", createQuery)
	return (*conn.inner).Exec(s.ctx, createQuery)
}
