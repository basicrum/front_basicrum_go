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
