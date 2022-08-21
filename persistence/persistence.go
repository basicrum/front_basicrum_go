package persistence

import (
	"context"
	"errors"
	"log"
	"net/http"

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
}

func New(s server, a auth, opts *opts) (*persistence, error) {
	if conn := s.open(&a); conn != nil {
		// @TODO: Move uaP dependency outside the persistance
		// We need to ge the Regexes from here: https://github.com/ua-parser/uap-core/blob/master/regexes.yaml
		uaP, err := uaparser.New("./assets/uaparser_regexes.yaml")
		if err != nil {
			log.Fatal(err)
		}

		return &persistence{s, connection{conn, a}, uaP, opts}, nil
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

func (p *persistence) Save(req *http.Request, table string) {
	tPrefix := &p.opts.prefix

	if table != "" {
		table = *tPrefix + "_" + table
	}

	p.server.save(&p.conn, eventPayload(req, p.uaP), table)
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
