package dao

import (
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// New Clickhouse connection
// nolint: revive
func NewConnection(s server, a auth) (driver.Conn, error) {
	return clickhouse.Open(&clickhouse.Options{
		Addr: []string{s.addr()},
		Auth: clickhouse.Auth{
			Database: s.db,
			Username: a.user,
			Password: a.pwd,
		},
		Debug:           false,
		ConnMaxLifetime: time.Hour,
	})
}
