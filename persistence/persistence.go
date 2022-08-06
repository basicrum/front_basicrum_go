package persistence

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
	Events chan *event
	opts   *opts
}

func New(s server, a auth, opts *opts) (*persistence, error) {
	if conn := s.open(&a); conn != nil {
		return &persistence{s, connection{conn, a}, make(chan *event), opts}, nil
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

func (p *persistence) Run() {
	for {
		select {
		case event := <-p.Events:
			if event != nil {
				go p.server.save(&p.conn, event.payload(), event.name)
			}
		}
	}
}

// START - Used for integration tests. Keeping ti dirty for now.
// @todo: Refactor or move big part of this to testing utility class.

func (p *persistence) RecycleTables() {
	p.server.RecycleTables(&p.conn)
}

func (p *persistence) CountRecords() uint64 {
	return p.server.countRecords(&p.conn)
}

// END - Used for integration tests.
