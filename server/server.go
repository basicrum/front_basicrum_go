package server

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/basicrum/front_basicrum_go/backup"
	"github.com/basicrum/front_basicrum_go/service"
	"github.com/rs/cors"
)

// Server represents http or https server
type Server struct {
	port           string
	service        *service.Service
	backup         backup.IBackup
	certFile       string
	keyFile        string
	tlsConfig      *tls.Config
	handlerAdapter func(http.Handler) http.Handler
	server         *http.Server
	listener       net.Listener
}

// WithHandlerAdapter creates server with handler wrapper function (used by Let's encrypt certificate manager)
func WithHandlerAdapter(handlerAdapter func(http.Handler) http.Handler) func(*Server) {
	return func(s *Server) {
		s.handlerAdapter = handlerAdapter
	}
}

// WithTLSConfig creates server with TLS configuration
func WithTLSConfig(tlsConfig *tls.Config) func(*Server) {
	return func(s *Server) {
		s.tlsConfig = tlsConfig
	}
}

// WithSSLFile creates server with SSL and certificate/key files
func WithSSLFile(certFile, keyFile string) func(*Server) {
	return func(s *Server) {
		s.certFile = certFile
		s.keyFile = keyFile
	}
}

// WithListener creates server with custom listener
func WithListener(listener net.Listener) func(*Server) {
	return func(s *Server) {
		s.listener = listener
	}
}

// WithPort creates server with port
func WithPort(port string) func(*Server) {
	return func(s *Server) {
		s.port = port
	}
}

// New creates a new http or https server
func New(
	processService *service.Service,
	backupService backup.IBackup,
	options ...func(*Server),
) *Server {
	result := &Server{
		service: processService,
		backup:  backupService,
		handlerAdapter: func(h http.Handler) http.Handler {
			return h
		},
	}
	for _, o := range options {
		o(result)
	}
	return result
}

// Serve start http or https server and blocks
func (s *Server) Serve() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/beacon/catcher", s.catcher)
	mux.HandleFunc("/health", s.health)
	handler := cors.Default().Handler(mux)
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.handlerAdapter(handler),
		// https://deepsource.io/directory/analyzers/go/issues/GO-S2114
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	if s.listener != nil {
		log.Println("starting server with listener")
		return s.server.Serve(s.listener)
	}
	if s.certFile != "" || s.keyFile != "" {
		log.Printf("starting https server on port[%v] with certFile[%v] keyFile[%v]", s.port, s.certFile, s.keyFile)
		return s.server.ListenAndServeTLS(s.certFile, s.keyFile)
	}
	log.Printf("starting http server on port[%v]", s.port)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shutdowns the http server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return errors.New("server is not started")
	}
	return s.server.Shutdown(ctx)
}
