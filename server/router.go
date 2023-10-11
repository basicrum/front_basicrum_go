package server

import "net/http"

func (s *Server) setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/beacon/catcher", s.catcher)
	mux.HandleFunc("/private/hostnames", s.hostnames)
	mux.HandleFunc("/health", s.health)
}
