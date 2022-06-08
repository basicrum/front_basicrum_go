package beacon

import (
	"fmt"
	"log"
	"testing"

	"github.com/ua-parser/uap-go/uaparser"
)

func TestBasic(t *testing.T) {
	b := Beacon{
		Pt_Lcp:     "230",
		U:          "https//:www.example.com/url",
		Nt_Con_End: "1653989622106",
		Nt_Con_St:  "1653989622032",
	}

	// We need to ge the Regexes from here: https://github.com/ua-parser/uap-core/blob/master/regexes.yaml
	uaP, err := uaparser.New("../assets/uaparser_regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	rE := ConvertToRumEvent(b, uaP)

	if rE.Connect_Duration != "74" {
		t.Errorf("Error")
	}

	//	fmt.Println(rE)

	fmt.Println(rE.Cumulative_Layout_Shift)

	if rE.Cumulative_Layout_Shift != "" {
		t.Errorf("Error")
	}
}
