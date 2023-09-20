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
func do(params []any, backupRootDir string) {
	backupsList := makeBackupList(params)
	saveBackupList(backupsList, backupRootDir)
}

// nolint: revive
func makeBackupList(params []any) map[string]string {
	backupsList := make(map[string]string)
	for _, p := range params {
		backupItem(backupsList, p)
	}
	return backupsList
}

func backupItem(backupsList map[string]string, p any) {
	v, ok := p.(url.Values)
	if !ok {
		// Can't assert, handle error.
		return
	}
	appendLine(backupsList, makeKey(v), makeValue(v))
}

func makeKey(v url.Values) string {
	urlValue, parseErr := url.Parse(v.Get("u"))
	if parseErr != nil {
		log.Print(parseErr)
	}
	return strings.ReplaceAll(urlValue.Hostname(), ".", "_")
}

func makeValue(v url.Values) string {
	flatten := flattenMap(v)

	dataJson, reqDataErr := json.Marshal(flatten)
	if reqDataErr != nil {
		log.Print(reqDataErr)
	}
	return string(dataJson)
}

func appendLine(m map[string]string, key, value string) {
	if _, keyFound := m[key]; !keyFound {
		m[key] = ""
	}
	m[key] += value + "\n"
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
	for host, lines := range backupsList {
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

		if _, err = f.WriteString(lines); err != nil {
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
