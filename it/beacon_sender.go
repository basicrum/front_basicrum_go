package it

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type BeaconSender struct {
	sender *HttpSender
}

func newBeaconSender(
	httpSender *HttpSender,
) *BeaconSender {
	return &BeaconSender{
		httpSender,
	}
}

func (b *BeaconSender) Send(path string) {
	requests, err := b.readFiles(path)
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
	for i, v := range requests {
		err := b.httpPost(v)
		if err != nil {
			log.Fatalf("Request[%d]: %v failed: %v", i, v, err)
		}
	}
}

func (b *BeaconSender) readFiles(path string) ([]url.Values, error) {
	files, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	result := []url.Values{}
	for _, file := range files {
		items, err := b.scanFile(file)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
	}
	return result, nil
}

func (b *BeaconSender) scanFile(fileName string) ([]url.Values, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to read file[%v]: %w", fileName, err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var result []url.Values
	for scanner.Scan() {
		var beaconData map[string]string
		err = json.Unmarshal([]byte(scanner.Text()), &beaconData)
		if err != nil {
			return nil, fmt.Errorf("unable to parse json[%s]: %w", []byte(scanner.Text()), err)
		}

		result = append(result, b.makeUrlValues(beaconData))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan file[%v]: %w", fileName, err)
	}
	return result, nil
}

func (b *BeaconSender) makeUrlValues(beaconDataMap map[string]string) url.Values {
	result := url.Values{}
	for k, v := range beaconDataMap {
		result.Set(k, v)
	}
	return result
}

func (b *BeaconSender) httpPost(params url.Values) error {
	headers, err := b.parseHeaders(params)
	if err != nil {
		return fmt.Errorf("parse headers error: %w", err)
	}

	req, err := http.NewRequest("POST", b.sender.BuildUrl("/beacon/catcher"), strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("newRequest error: %w", err)
	}

	req.Header.Add("User-Agent", makeUserAgent(headers))
	req.Header.Add("Cf-Ipcountry", makeCountryCode(headers))
	req.Header.Add("Cf-Ipcity", makeCityName(headers))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	b.sender.Send(req, http.StatusNoContent, "")
	return nil
}

func (b *BeaconSender) parseHeaders(params url.Values) (map[string][]string, error) {
	var headers map[string][]string

	err := json.Unmarshal([]byte(params.Get("request_headers")), &headers)
	if err != nil {
		return nil, fmt.Errorf("unable to parse json[%s]: %w", params, err)
	}

	return headers, nil
}

func makeUserAgent(headers map[string][]string) string {
	if userAgent, ok := headers["User-Agent"]; ok {
		return userAgent[0]
	}
	return ""
}

func makeCountryCode(headers map[string][]string) string {
	return headers["Cf-Ipcountry"][0]
}

func makeCityName(headers map[string][]string) string {
	return headers["Cf-Ipcity"][0]
}
