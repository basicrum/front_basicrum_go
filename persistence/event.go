package persistence

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/ua-parser/uap-go/uaparser"
)

// TODO !!! beacon logic must reside in beacon pkg, for now it is just a copy-paste from main
func eventPayload(req *http.Request, uaP *uaparser.Parser) string {
	b := beacon.FromRequestParams(&req.Form, req.UserAgent(), req.Header)
	re := beacon.ConvertToRumEvent(b, uaP)
	jsonValue, err := json.Marshal(re)

	if err != nil {
		log.Printf("json parsing error: %+v", err)
		return ""
	}

	return string(jsonValue)
}
