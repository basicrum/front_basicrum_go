package persistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func ConnectClickHouse(host string, port string, dbName string, userName string, password string) (error, driver.Conn) {
	addr := host + ":" + port

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: dbName,
			Username: userName,
			Password: password,
		},
		Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})

	return err, conn
}

func SaveInClickHouse(ctx context.Context, conn driver.Conn, table string, rumEvent string) {

	inertQuery := fmt.Sprintf(
		`INSERT INTO %s SETTINGS input_format_skip_unknown_fields = true FORMAT JSONEachRow
			%s`, table, rumEvent)

	insErr := conn.AsyncInsert(ctx, inertQuery, false)

	if insErr != nil {
		log.Fatal(insErr)
	}

}

func RecycleTables(ctx context.Context, conn driver.Conn) {

	dropQuery := `DROP TABLE IF EXISTS integration_test_webperf_rum_events`

	dropErr := conn.Exec(ctx, dropQuery)

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
	
		country_code                    FixedString(2)
	)
		ENGINE = MergeTree()
		PARTITION BY toYYYYMMDD(event_date)
		ORDER BY (device_type, event_date)
		SETTINGS index_granularity = 8192`

	createErr := conn.Exec(ctx, createQuery)

	if createErr != nil {
		log.Fatal(createErr)
	}
}

func CountRecords(ctx context.Context, conn driver.Conn) {
	rows, err := conn.Query(ctx, "SELECT count(*) FROM integration_test_webperf_rum_events")
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
