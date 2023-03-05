package dao

import (
	"strconv"
)

type server struct {
	host string
	port int16
	db   string
}

func (s *server) addr() string {
	return s.host + ":" + strconv.FormatInt(int64(s.port), 10)
}

type auth struct {
	user string
	pwd  string
}

type opts struct {
	prefix string
}

// Server creates the datastore (click house) options
// nolint: revive
func Server(host string, port int16, db string) server {
	return server{host, port, db}
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
