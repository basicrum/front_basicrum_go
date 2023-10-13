package it

import (
	"context"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

const baseTableName = "webperf_rum_events"

type opts struct {
	prefix string
}

type IntegrationDao struct {
	conn driver.Conn
	opts opts
}

func NewIntegrationDao(conn driver.Conn, opts opts) *IntegrationDao {
	return &IntegrationDao{
		conn,
		opts,
	}
}

func Opts(prefix string) opts {
	return opts{prefix}
}

func (p *IntegrationDao) fullTableName(name string) string {
	return p.opts.prefix + name
}

func (p *IntegrationDao) RecycleTables() {
	dropQuery := fmt.Sprintf("TRUNCATE TABLE %v", p.fullTableName(baseTableName))
	dropErr := p.conn.Exec(context.Background(), dropQuery)
	if dropErr != nil {
		log.Print(dropErr)
	}
}

func (p *IntegrationDao) CountRecords(criteria string) int {
	query := fmt.Sprintf("SELECT count(*) FROM %v %v", p.fullTableName(baseTableName), criteria)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if !rows.Next() {
		return 0
	}

	var result uint64
	if err := rows.Scan(&result); err != nil {
		panic(err)
	}

	return int(result)
}
