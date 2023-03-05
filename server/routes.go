package server

import (
	"log"
	"net/http"
	"time"

	"github.com/basicrum/front_basicrum_go/types"
)

func (s *Server) catcher(w http.ResponseWriter, r *http.Request) {
	// return no cache headers
	s.headersNoCache(w)

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

func newEventFromRequest(r *http.Request) (*types.Event, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	form := r.Form
	// We need this in case we would like to re-import beacons
	// Also created_at is used for event date when we persist data in the DB
	if !form.Has("created_at") {
		form.Set("created_at", time.Now().UTC().Format("2006-01-02 15:04:05"))
	}
	return types.NewEvent(form, r.Header, r.UserAgent()), nil
}

func (*Server) headersNoCache(w http.ResponseWriter) {
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
	w.WriteHeader(http.StatusNoContent)
}

func (*Server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte("ok"))
}
