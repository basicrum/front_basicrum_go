package server

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/basicrum/front_basicrum_go/types"
)

func (s *Server) catcher(w http.ResponseWriter, r *http.Request) {
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

	// Archiving logic - save the event to a file
	s.backup.SaveAsync(event)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
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

func (s *Server) responseNoContent(w http.ResponseWriter) {
	s.headersNoCache(w, http.StatusNoContent)
}

func (s *Server) responseOK(w http.ResponseWriter) {
	s.headersNoCache(w, http.StatusOK)
	_, _ = w.Write([]byte("ok"))
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
