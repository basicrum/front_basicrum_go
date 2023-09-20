package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
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
	appendLine(backupsList, makeKeyHostname(v), makeValue(v))
}

func makeKeyHostname(v url.Values) string {
	urlValue, parseErr := url.Parse(v.Get("u"))
	if parseErr != nil {
		log.Print(parseErr)
	}
	return strings.ReplaceAll(urlValue.Hostname(), ".", "_")
}

func makeValue(v url.Values) string {
	return toJSON(flattenMap(v))
}

func toJSON(flatten any) string {
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
	for host, lines := range backupsList {
		saveLinesToFile(backupRootDir, host, lines)
	}
}

func saveLinesToFile(backupRootDir string, host string, lines string) {
	filename := makeFilePath(backupRootDir, host)
	err := appendToFile(filename, lines)
	if err != nil {
		log.Print(err)
	}
}

func appendToFile(filename string, lines string) error {
	f := openOrCreateFileForAppend(filename)
	go func() {
		if err := f.Close(); err != nil {
			log.Print(err)
		}
	}()

	_, err := f.WriteString(lines)
	return err
}

func writeToFile(filename string, lines string, factory CompressionWriterFactory) error {
	f, err := openFile(filename)
	if err != nil {
		return err
	}
	go func() {
		if err := f.Close(); err != nil {
			log.Print(err)
		}
	}()
	compressionWriter, err := factory.Create(f)
	if err != nil {
		return err
	}
	go func() {
		if err := compressionWriter.Close(); err != nil {
			log.Print(err)
		}
	}()

	return writeString(compressionWriter, lines)
}

func writeString(f io.Writer, s string) error {
	_, err := f.Write([]byte(s))
	return err
}

func openOrCreateFileForAppend(filename string) *os.File {
	err := os.MkdirAll(path.Dir(filename), os.ModePerm)
	if err != nil {
		log.Print(err)
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Print(err)
	}
	return f
}

func openFile(filename string) (*os.File, error) {
	err := os.MkdirAll(path.Dir(filename), os.ModePerm)
	if err != nil {
		return nil, err
	}

	err = os.Remove(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func makeFilePath(backupRootDir string, host string) string {
	now := time.Now()
	return makeHourPath(
		makeDayPath(backupRootDir, host, now),
		now,
	)
}

func makeDayPath(backupRootDir string, host string, nowUTC time.Time) string {
	return filepath.Join(backupRootDir, host, dateUTC(nowUTC))
}

func makeArchiveDayPath(backupRootDir, host string, now time.Time) string {
	return filepath.Join(backupRootDir, host, dateUTC(now)+".json.lines")
}

func makeArchiveDayMetaPath(backupRootDir, host string, now time.Time) string {
	return filepath.Join(backupRootDir, host, dateUTC(now)+".meta.lines")
}

func makeHourPath(parent string, now time.Time) string {
	return filepath.Join(parent, hourUTC(now)+".json.lines")
}

func dateUTC(now time.Time) string {
	nowUTC := now.UTC()
	return fmt.Sprintf("%v-%v-%v", nowUTC.Year(), int(nowUTC.Month()), nowUTC.Day())
}

func hourUTC(now time.Time) string {
	nowUTC := now.UTC()
	return strconv.Itoa(nowUTC.Hour())
}
