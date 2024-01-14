package server

import "net/http"

func (s *Server) setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/beacon/catcher", s.catcher)
	mux.HandleFunc("/health", s.health)
}
