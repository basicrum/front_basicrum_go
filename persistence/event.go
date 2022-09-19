package persistence

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/ua-parser/uap-go/uaparser"
)

type event struct {
	name      string
	reqParams *url.Values
	headers   *http.Header
	userAgent string
}

func (p *persistence) Event(reqParams *url.Values, headers *http.Header, userAgent string) *event {
	return &event{"webperf_rum_events", reqParams, headers, userAgent}
}

// TODO !!! beacon logic must reside in beacon pkg, for now it is just a copy-paste from main
func (e *event) payload(uaP *uaparser.Parser) string {
	b := beacon.FromRequestParams(e.reqParams, e.userAgent, e.headers)
	re := beacon.ConvertToRumEvent(b, uaP)
	jsonValue, err := json.Marshal(re)

	if err != nil {
		log.Printf("json parsing error: %+v", err)
		return ""
	}

	return string(jsonValue)
}
