package dao

import (
	"context"
	"database/sql"
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

func (p *DAO) GetSubscriptions() (map[string]types.Subscription, error) {
	columns := "subscription_id, subscription_expire_at"
	query := fmt.Sprintf("SELECT %v FROM %v%v", columns, p.prefix, baseOwnerHostsTableName)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("get subscriptions failed: %w", err)
	}
	// defer rows.Close()

	subscriptions := make(map[string]types.Subscription)
	for rows.Next() {
		var subscription types.Subscription
		if err := rows.Scan(&subscription.ID, &subscription.ExpiresAt); err != nil {
			return subscriptions, err
		}
		subscriptions[subscription.ID] = subscription
	}

	if err = rows.Err(); err != nil {
		return subscriptions, err
	}
	fmt.Println(subscriptions)
	return subscriptions, nil
}

func (p *DAO) GetSubscription(id string) (types.Subscription, error) {
	var subscription types.Subscription

	columns := "subscription_id, subscription_expire_at"
	whereClause := "WHERE hostname='" + id + "'"
	query := fmt.Sprintf("SELECT %v FROM %v%v %v", columns, p.prefix, baseOwnerHostsTableName, whereClause)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return subscription, fmt.Errorf("get subscription failed: %w", err)
	}
	// defer rows.Close()

	if !rows.Next() {
		return subscription, nil
	}
	err = rows.Scan(&subscription.ID, &subscription.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return subscription, fmt.Errorf("subscription with id: %s not found", id)
		}
		return subscription, fmt.Errorf("get subscription failed: %w", err)
	}

	return subscription, nil
}
