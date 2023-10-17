package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/types"
)

const (
	grafanaUserHeader  = "X-Grafana-User"
	privateTokenHeader = "X-Token"
)

func (s *Server) catcher(w http.ResponseWriter, r *http.Request) {
	sConf, err := config.GetStartupConfig()
	if err != nil {
		log.Fatal(err)
	}
	// return no cache headers
	s.responseNoContent(w)

	// create an event from http request
	event, err := newEventFromRequest(r)
	if err != nil {
		log.Printf("failed to parse request %+v", err)
		return
	}

	// Persist Event async in ClickHouse
	s.service.SaveAsync(event)

	if sConf.Backup.Enabled {
		// Archiving logic - save the event to a file
		s.backup.SaveAsync(event)
	}
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	s.responseOK(w)
}

func (s *Server) hostnames(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get(privateTokenHeader)
	if token != s.privateAPIToken {
		s.responseError(w, fmt.Errorf("wrong header[%v]", privateTokenHeader))
		return
	}
	switch r.Method {
	case http.MethodPost:
		s.createHostname(w, r)
	case http.MethodDelete:
		s.deleteHostname(w, r)
	default:
		s.responseError(w, fmt.Errorf("unsupported method[%v]", r.Method))
	}
}

// CreateHostnameRequest is create hostname request
type CreateHostnameRequest struct {
	Hostname string `json:"hostname"`
}

// Validate create hostname request
func (r *CreateHostnameRequest) Validate() error {
	if r.Hostname == "" {
		return errors.New("hostname is required")
	}
	return nil
}

func (s *Server) createHostname(w http.ResponseWriter, r *http.Request) {
	var request CreateHostnameRequest
	if err := s.parseRequest(r, &request); err != nil {
		s.responseError(w, err)
		return
	}
	username, err := s.parseGrafanaUser(r)
	if err != nil {
		s.responseError(w, err)
		return
	}
	// nolint: contextcheck
	err = s.service.RegisterHostname(request.Hostname, username)
	if err != nil {
		s.responseError(w, err)
		return
	}
	s.responseOK(w)
}

// DeleteHostnameRequest is delete hostname request
type DeleteHostnameRequest struct {
	Hostname string `json:"hostname"`
}

// Validate delete hostname request
func (r *DeleteHostnameRequest) Validate() error {
	if r.Hostname == "" {
		return errors.New("hostname is required")
	}
	return nil
}

func (s *Server) deleteHostname(w http.ResponseWriter, r *http.Request) {
	var request DeleteHostnameRequest
	if err := s.parseRequest(r, &request); err != nil {
		s.responseError(w, err)
		return
	}
	username, err := s.parseGrafanaUser(r)
	if err != nil {
		s.responseError(w, err)
		return
	}
	// nolint: contextcheck
	err = s.service.DeleteHostname(request.Hostname, username)
	if err != nil {
		s.responseError(w, err)
		return
	}
	s.responseOK(w)
}

func newEventFromRequest(r *http.Request) (*types.Event, error) {
	form, err := parseEventForm(r)
	if err != nil {
		return nil, err
	}
	ip := getIP(r)
	return types.NewEvent(form, r.Header, r.UserAgent(), ip), nil
}

func parseEventForm(r *http.Request) (url.Values, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	applyEventFormDefaults(r.Form)
	return r.Form, nil
}

func applyEventFormDefaults(form url.Values) {
	// We need this in case we would like to re-import beacons
	// Also created_at is used for event date when we persist data in the DB
	if !form.Has("created_at") {
		form.Set("created_at", time.Now().UTC().Format("2006-01-02 15:04:05"))
	}
}

func getIP(r *http.Request) string {
	var result string
	var temp string

	result = r.RemoteAddr
	temp = r.Header.Get("X-Forwarded-For")
	if result != "" {
		result = temp
	}

	parts := strings.Split(result, ",")
	return parts[0]
}

func (*Server) parseRequest(r *http.Request, request Validator) error {
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}
	return request.Validate()
}

func (*Server) parseGrafanaUser(r *http.Request) (string, error) {
	result := r.Header.Get(grafanaUserHeader)
	if result == "" {
		return "", fmt.Errorf("required header[%v]", grafanaUserHeader)
	}
	return result, nil
}

func (s *Server) responseNoContent(w http.ResponseWriter) {
	s.headersNoCache(w, http.StatusNoContent)
}

func (s *Server) responseOK(w http.ResponseWriter) {
	s.headersNoCache(w, http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) responseError(w http.ResponseWriter, err error) {
	s.headersNoCache(w, http.StatusBadRequest)
	_, _ = w.Write([]byte(err.Error()))
}

func (*Server) headersNoCache(w http.ResponseWriter, statusCode int) {
	// @todo: Check if we need to add more response headers
	// access-control-allow-credentials: true
	// access-control-allow-origin: *
	// cache-control: no-cache, no-store, must-revalidate
	// content-length: 0
	// content-type: text/plain
	// cross-origin-resource-policy: cross-origin
	// date: Sat, 25 Jun 2022 10:40:18 GMT
	// expires: Fri, 01 Jan 1990 00:00:00 GMT
	// pragma: no-cache
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
	w.WriteHeader(statusCode)
}
