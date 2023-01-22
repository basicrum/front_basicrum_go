package persistence

import (
	"context"
	"errors"
	"log"

	"github.com/ua-parser/uap-go/uaparser"
)

const baseTableName = "webperf_rum_events"

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
	server    server
	conn      connection
	uaP       *uaparser.Parser
	opts      *opts
	Events    chan *event
	tableName string
}

// New creates persistance service
// nolint: revive
func New(s server, a auth, opts *opts, uaP *uaparser.Parser) (*persistence, error) {
	if conn := s.open(&a); conn != nil {
		tableName := opts.prefix + baseTableName
		return &persistence{s, connection{conn, a}, uaP, opts, make(chan *event), tableName}, nil
	}

	return nil, errors.New("connection to the server failed")
}

// Server creates the datastore (click house) options
// nolint: revive
func Server(host string, port int16, db string) server {
	return server{host, port, db, context.Background()}
}

// Auth creates the authentication options for persistance service
// nolint: revive
func Auth(user string, pwd string) auth {
	return auth{user, pwd}
}

// Opts creates the options for persistance service
// nolint: revive
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
		go p.server.save(&p.conn, event.payload(p.uaP), p.tableName)
	}
}

// CreateTable creates the table if not exists
func (p *persistence) CreateTable() error {
	tableExist, err := p.server.CheckTableExist(&p.conn, p.tableName)
	if err != nil {
		return err
	}
	if tableExist {
		log.Printf("table already exists")
		return nil
	}
	return p.server.CreateTable(&p.conn, p.tableName)
}
