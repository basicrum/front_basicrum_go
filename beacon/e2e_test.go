package beacon

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ua-parser/uap-go/uaparser"
)

func TestE2E(t *testing.T) {
	// We need to ge the Regexes from here: https://github.com/ua-parser/uap-core/blob/master/regexes.yaml
	uaP, err := uaparser.New("./../assets/uaparser_regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	r := RequestFromFile("./data/req-0")

	b := FromRequestParams(&r.Form, r.UserAgent(), r.Header)

	re := ConvertToRumEvent(b, uaP)

	jsonValue, _ := json.Marshal(re)

	fmt.Println(string(jsonValue))
}

func RequestFromFile(fName string) *http.Request {

	file, err := os.Open(fName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	nextIsParams := false
	paramsStr := ""

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		if nextIsParams == true {
			paramsStr = scanner.Text()
			break
		}

		if scanner.Text() == "" {
			nextIsParams = true
		}
		// fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "localhost:80", bytes.NewBuffer([]byte("")))

	if err != nil {
		log.Fatal(err)
	}

	req.URL.RawQuery = paramsStr

	req.ParseForm()

	return req
}
