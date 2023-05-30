package beacon

import (
	"log"
	"testing"

	"github.com/basicrum/front_basicrum_go/types"
	"github.com/ua-parser/uap-go/uaparser"
)

func TestBasic(t *testing.T) {
	b := Beacon{
		Pt_Lcp:     "230",
		U:          "https//:www.example.com/url",
		Nt_Con_End: "1653989622106",
		Nt_Con_St:  "1653989622032",
	}
	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.82 Safari/537.36"

	// We need to ge the Regexes from here: https://github.com/ua-parser/uap-core/blob/master/regexes.yaml
	uaP, err := uaparser.New("../assets/uaparser_regexes.yaml")
	if err != nil {
		panic(err)
	}

	event := &types.Event{UserAgent: userAgent}
	rE := ConvertToRumEvent(b, event, uaP, nil)

	if rE.Connect_Duration != "74" {
		t.Errorf("Error")
	}

	log.Printf("rum event[%+v]", rE)

	log.Printf("rum event Cumulative_Layout_Shift[%v]", rE.Cumulative_Layout_Shift)

	if rE.Cumulative_Layout_Shift != "" {
		t.Errorf("Error")
	}
}
