package types

import (
	"net/http"
	"net/url"
)

// Event contains the http request of catcher - body/query parameters, headers and user agent
type Event struct {
	RequestParameters url.Values
	Headers           http.Header
	UserAgent         string
}

// NewEvent creates a new event to be stored
func NewEvent(requestParameters url.Values, headers http.Header, userAgent string) *Event {
	return &Event{requestParameters, headers, userAgent}
}
