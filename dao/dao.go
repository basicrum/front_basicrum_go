package dao

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

const baseTableName = "webperf_rum_events"

// DAO is data access object for clickhouse database
type DAO struct {
	conn  clickhouse.Conn
	table string
}

// New creates persistance service
// nolint: revive
func New(s server, a auth, opts *opts) (*DAO, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{s.addr()},
		Auth: clickhouse.Auth{
			Database: s.db,
			Username: a.user,
			Password: a.pwd,
		},
		Debug:           false,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("clickhouse connection failed: %w", err)
	}
	table := opts.prefix + baseTableName
	return &DAO{
		conn:  conn,
		table: table,
	}, nil
}

// Save stores data into table in clickhouse database
func (p *DAO) Save(data string) error {
	if data == "" {
		return fmt.Errorf("clickhouse invalid data for table %s: %s", p.table, data)
	}
	query := fmt.Sprintf(
		"INSERT INTO %s SETTINGS input_format_skip_unknown_fields = true FORMAT JSONEachRow %s",
		p.table,
		data,
	)
	err := p.conn.AsyncInsert(context.Background(), query, false)
	if err != nil {
		return fmt.Errorf("clickhouse insert failed: %w", err)
	}
	return nil
}

// CreateTableIfNotExist creates the table if not exists
func (p *DAO) CreateTableIfNotExist() error {
	tableExist, err := p.CheckTableExist()
	if err != nil {
		return err
	}
	if tableExist {
		log.Printf("table already exists")
		return nil
	}
	return p.CreateTable()
}

// CheckTableExist checks if table exists
func (p *DAO) CheckTableExist() (bool, error) {
	query := fmt.Sprintf(`EXISTS %s`, p.table)
	rows, err := p.conn.Query(context.Background(), query)
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
func (p *DAO) CreateTable() error {
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
		dom_link_css                    Nullable(UInt16),

		mob_etype                       LowCardinality(Nullable(String)),
		mob_dl                          Nullable(UInt16),
		mob_rtt                         Nullable(UInt16)

	)
		ENGINE = MergeTree()
		PARTITION BY toYYYYMMDD(event_date)
		ORDER BY (hostname, event_date)
		SETTINGS index_granularity = 8192`, p.table)

	log.Printf("creating table with query: %v", createQuery)
	return p.conn.Exec(context.Background(), createQuery)
}
