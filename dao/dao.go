package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"

	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/basicrum/front_basicrum_go/types"
)

const (
	baseTableName           = "webperf_rum_events"
	baseHostsTableName      = "webperf_rum_hostnames"
	baseOwnerHostsTableName = "webperf_rum_own_hostnames"
	tablePrefixPlaceholder  = "{prefix}"
	bufferSize              = 1024
)

// DAO is data access object for clickhouse database
type DAO struct {
	conn   clickhouse.Conn
	table  string
	prefix string
}

// New creates persistance service
// nolint: revive
func New(conn clickhouse.Conn, opts *opts) *DAO {
	return &DAO{
		conn:   conn,
		table:  fullTableName(opts),
		prefix: opts.prefix,
	}
}

func fullTableName(opts *opts) string {
	return opts.prefix + baseTableName
}

// Close the clickhouse connection
func (p *DAO) Close() error {
	return p.conn.Close()
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

// SaveHost stores hostname data into table in clickhouse database
func (p *DAO) SaveHost(event beacon.HostnameEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(
		"INSERT INTO %s%s SETTINGS input_format_skip_unknown_fields = true FORMAT JSONEachRow %s",
		p.prefix,
		baseHostsTableName,
		data,
	)
	err = p.conn.AsyncInsert(context.Background(), query, false)
	if err != nil {
		return fmt.Errorf("clickhouse insert failed: %w", err)
	}
	return nil
}

// InsertOwnerHostname inserts a new hostname
func (p *DAO) InsertOwnerHostname(item types.OwnerHostname) error {
	query := fmt.Sprintf(
		"INSERT INTO %s%s(username, hostname, subscription_id, subscription_expire_at) VALUES(?,?,?,?)",
		p.prefix,
		baseOwnerHostsTableName,
	)
	return p.conn.Exec(context.Background(), query, item.Username, item.Hostname, item.Subscription.ID, item.Subscription.ExpiresAt)
}

// DeleteOwnerHostname deletes the hostname
func (p *DAO) DeleteOwnerHostname(hostname, username string) error {
	query := fmt.Sprintf(
		"DELETE FROM %s%s WHERE hostname = ? AND username = ?",
		p.prefix,
		baseOwnerHostsTableName,
	)
	return p.conn.Exec(context.Background(), query, hostname, username)
}
