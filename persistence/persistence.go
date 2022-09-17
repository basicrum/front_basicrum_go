package persistence

import (
	"context"
	"errors"

	"github.com/ua-parser/uap-go/uaparser"
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
	uaP    *uaparser.Parser
	opts   *opts
	Events chan *event
}

func New(s server, a auth, opts *opts, uaP *uaparser.Parser) (*persistence, error) {
	if conn := s.open(&a); conn != nil {
		return &persistence{s, connection{conn, a}, uaP, opts, make(chan *event)}, nil
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
				tPrefix := &p.opts.prefix

				table := event.name

				if *tPrefix != "" {
					table = *tPrefix + "_" + table
				}

				go p.server.save(&p.conn, event.payload(p.uaP), table)
			}
		}
	}
}
