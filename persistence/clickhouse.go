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
		// Debug:           true,
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
			%s`,
		table,
		rumEvent)

	insErr := conn.AsyncInsert(ctx, inertQuery, false)

	if insErr != nil {
		log.Fatal(insErr)
	}

}
