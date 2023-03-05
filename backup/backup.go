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

// do saves parameters to a file in backup directory
func do(params []interface{}, backupRootDir string) {
	backupsList := makeBackupList(params)
	saveBackupList(backupsList, backupRootDir)
}

// nolint: revive
func makeBackupList(params []interface{}) map[string]string {
	backupsList := make(map[string]string)
	for _, p := range params {
		v, ok := p.(url.Values)
		if !ok {
			// Can't assert, handle error.
			continue
		}

		flatten := flattenMap(v)

		dataJson, reqDataErr := json.Marshal(flatten)
		if reqDataErr != nil {
			log.Print(reqDataErr)
		}

		urlValue, parseErr := url.Parse(v.Get("u"))
		if parseErr != nil {
			log.Print(parseErr)
		}

		hostNormalized := strings.ReplaceAll(urlValue.Hostname(), ".", "_")
		if _, containsHost := backupsList[hostNormalized]; !containsHost {
			backupsList[hostNormalized] = ""
		}

		backupsList[hostNormalized] += string(dataJson) + "\n"
	}
	return backupsList
}

func flattenMap(v url.Values) map[string]string {
	flatten := make(map[string]string)
	for k := range v {
		flatten[k] = v.Get(k)
	}
	return flatten
}

func saveBackupList(backupsList map[string]string, backupRootDir string) {
	// Date path
	datePath := getDateUtcPath()
	utcHour := getUtcHour()
	for host, data := range backupsList {
		dirPath := backupRootDir + host + "/" + datePath

		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			log.Print(err)
		}

		filename := dirPath + "/" + utcHour + ".json.lines"
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			log.Print(err)
		}
		// nolint: revive
		defer f.Close()

		if _, err = f.WriteString(data); err != nil {
			log.Print(err)
		}
	}
}

func getDateUtcPath() string {
	t := time.Now().UTC()
	utcYear := strconv.Itoa(t.Year())
	utcMonth := strconv.Itoa(int(t.Month()))
	utcDay := strconv.Itoa(t.Day())

	return utcYear + "-" + utcMonth + "-" + utcDay
}

func getUtcHour() string {
	t := time.Now().UTC()
	return strconv.Itoa(t.Hour())
}
