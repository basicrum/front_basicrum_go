package it

import (
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

type RumEventRow struct {
	Url                     string `ch:"url"`
	Cumulative_Layout_Shift string `ch:"cumulative_layout_shift"`
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

func (s *server) RecycleTables(conn *connection) {
	dropQuery := `TRUNCATE TABLE integration_test_webperf_rum_events`
	dropErr := (*conn.inner).Exec(s.ctx, dropQuery)
	if dropErr != nil {
		log.Print(dropErr)
	}
}

func (s *server) countRecords(conn *connection, criteria string) uint64 {
	query := "SELECT count(*) FROM integration_test_webperf_rum_events " + criteria
	rows, err := (*conn.inner).Query(s.ctx, query)
	if err != nil {
		log.Fatal(err)
	}

	var cnt uint64 = 0

	for rows.Next() {
		var (
			col1 uint64
		)

		if err := rows.Scan(&col1); err != nil {
			log.Fatal(err)
		}

		cnt = col1
	}
	rows.Close()

	return cnt
}

func (s *server) getFirstRow(conn *connection) RumEventRow {

	result := []RumEventRow{}

	err := (*conn.inner).Select(s.ctx, &result, "SELECT url, cumulative_layout_shift FROM integration_test_webperf_rum_events")

	if err != nil {
		log.Fatal(err)
	}

	return result[0]
}
