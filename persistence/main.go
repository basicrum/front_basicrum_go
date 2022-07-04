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

type persistence struct {
	server server
	conn   connection
}

func New(s server, a auth) (*persistence, error) {
	if conn := s.open(&a); conn != nil {
		return &persistence{s, connection{conn, a}}, nil
	}
	return nil, errors.New("connection to the server failed")
}

func Server(host string, port int16, db string) server {
	return server{host, port, db, context.Background()}
}

func Auth(user string, pwd string) auth {
	return auth{user, pwd}
}

func (p *persistence) Save(data []byte, name string) {
	p.server.save(&p.conn, string(data), name)
}
