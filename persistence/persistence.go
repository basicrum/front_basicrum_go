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

// New creates persistance service
func New(s server, a auth, opts *opts, uaP *uaparser.Parser) (*persistence, error) {
	if conn := s.open(&a); conn != nil {
		return &persistence{s, connection{conn, a}, uaP, opts, make(chan *event)}, nil
	}

	return nil, errors.New("connection to the server failed")
}

// Server creates the datastore (click house) options
func Server(host string, port int16, db string) server {
	return server{host, port, db, context.Background()}
}

// Auth creates the authentication options for persistance service
func Auth(user string, pwd string) auth {
	return auth{user, pwd}
}

// Opts creates the options for persistance service
func Opts(prefix string) *opts {
	return &opts{prefix}
}

// Run process the events from the channel and save them in datastore (click house)
func (p *persistence) Run() {
	for {
		event := <-p.Events
		if event == nil {
			continue
		}
		table := event.name

		tPrefix := p.opts.prefix
		if tPrefix != "" {
			table = tPrefix + "_" + table
		}

		go p.server.save(&p.conn, event.payload(p.uaP), table)
	}
}
