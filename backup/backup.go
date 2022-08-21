package backup

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func Do(params []interface{}, backupRootDir string) {

	backupsList := make(map[string]string)

	// Date path
	datePath := GetDateUtcPath()
	utcHour := GetUtcHour()

	for _, p := range params {

		v, ok := p.(url.Values)

		if !ok {
			// Can't assert, handle error.
			continue
		}

		flatten := make(map[string]string)

		for k := range v {
			flatten[k] = v.Get(k)
		}

		dataJson, reqDataErr := json.Marshal(flatten)

		if reqDataErr != nil {
			log.Print(reqDataErr)
		}

		url, parseErr := url.Parse(v.Get("u"))

		if parseErr != nil {
			log.Print(parseErr)
		}

		hostNormalized := strings.ReplaceAll(url.Hostname(), ".", "_")

		if _, containsHost := backupsList[hostNormalized]; !containsHost {
			backupsList[hostNormalized] = ""
		}

		backupsList[hostNormalized] = backupsList[hostNormalized] + string(dataJson) + "\n"
	}

	for host, data := range backupsList {
		dirPath := backupRootDir + host + "/" + datePath

		os.MkdirAll(dirPath, os.ModePerm)

		filename := dirPath + "/" + utcHour + ".json.lines"

		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)

		if err != nil {
			log.Print(err)
		}

		defer f.Close()

		if _, err = f.WriteString(data); err != nil {
			log.Print(err)
		}
	}
}

func GetDateUtcPath() string {
	t := time.Now().UTC()
	utcYear := strconv.Itoa(t.Year())
	utcMonth := strconv.Itoa(int(t.Month()))
	utcDay := strconv.Itoa(t.Day())

	return utcYear + "-" + utcMonth + "-" + utcDay
}

func GetUtcHour() string {
	t := time.Now().UTC()
	return strconv.Itoa(t.Hour())
}
