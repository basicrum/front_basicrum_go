package persistence

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/ua-parser/uap-go/uaparser"
)

type event struct {
	name string
	req  *http.Request
}

func (p *persistence) Event(r *http.Request) *event {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		return nil
	}

	name := "webperf_rum_events"
	if p.opts != nil {
		name = fmt.Sprintf("%s_%s", p.opts.prefix, name)
	}

	return &event{name, r}
}

// TODO !!! beacon logic must reside in beacon pkg, for now it is just a copy-paste from main
func (e *event) payload() string {
	// We need to ge the Regexes from here: https://github.com/ua-parser/uap-core/blob/master/regexes.yaml
	uaP, err := uaparser.New("./assets/uaparser_regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}
	b := beacon.FromRequestParams(&e.req.Form, e.req.UserAgent(), e.req.Header)
	re := beacon.ConvertToRumEvent(b, uaP)
	jsonValue, err := json.Marshal(re)
	if err != nil {
		log.Fatalf("json parsing error: %+v", err)
		return ""
	}
	return string(jsonValue)
}
