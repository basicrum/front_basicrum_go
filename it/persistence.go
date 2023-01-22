package it

import (
	"context"
	"errors"
)

const baseTableName = "webperf_rum_events"

type auth struct {
	user string
	pwd  string
}

type server struct {
	host      string
	port      int16
	db        string
	ctx       context.Context
	tableName string
}

type opts struct {
	prefix string
}

type Persistence struct {
	server server
	conn   connection
	opts   *opts
}

func New(s server, a auth, opts *opts) (*Persistence, error) {
	if conn := s.open(&a); conn != nil {
		return &Persistence{s, connection{conn, a}, opts}, nil
	}

	return nil, errors.New("connection to the server failed")
}

func Server(host string, port int16, db string, tablePrefix string) server {
	tableName := tablePrefix + baseTableName
	return server{host, port, db, context.Background(), tableName}
}

func Auth(user string, pwd string) auth {
	return auth{user, pwd}
}

func Opts(prefix string) *opts {
	return &opts{prefix}
}

func (p *Persistence) RecycleTables() {
	p.server.RecycleTables(&p.conn)
}

func (p *Persistence) CountRecords(criteria string) uint64 {
	return p.server.countRecords(&p.conn, criteria)
}

func (p *Persistence) GetFirstRow() RumEventRow {
	return p.server.getFirstRow(&p.conn)
}
