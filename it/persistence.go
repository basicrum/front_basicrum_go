package it

import (
	"context"
	"errors"
)

type auth struct {
	user string
	pwd  string
}

type server struct {
	host string
	port int16
	db   string
	ctx  context.Context
}

type opts struct {
	prefix string
}

type persistence struct {
	server server
	conn   connection
	opts   *opts
}

func New(s server, a auth, opts *opts) (*persistence, error) {
	if conn := s.open(&a); conn != nil {
		return &persistence{s, connection{conn, a}, opts}, nil
	}

	return nil, errors.New("connection to the server failed")
}

func Server(host string, port int16, db string) server {
	return server{host, port, db, context.Background()}
}

func Auth(user string, pwd string) auth {
	return auth{user, pwd}
}

func Opts(prefix string) *opts {
	return &opts{prefix}
}

func (p *persistence) RecycleTables() {
	p.server.RecycleTables(&p.conn)
}

func (p *persistence) CountRecords(criteria string) uint64 {
	return p.server.countRecords(&p.conn, criteria)
}

func (p *persistence) GetFirstRow() RumEventRow {
	return p.server.getFirstRow(&p.conn)
}
