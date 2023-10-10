package it

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type newStyleBeaconSender struct {
	client *http.Client
	host   string
	port   string
}

func newNewStyleBeaconSender(
	client *http.Client,
	host string,
	port string,
) *newStyleBeaconSender {
	return &newStyleBeaconSender{
		client: client,
		host:   host,
		port:   port,
	}
}

func (b *newStyleBeaconSender) Send(path string) {
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

func (b *newStyleBeaconSender) readFiles(path string) ([]url.Values, error) {
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

func (b *newStyleBeaconSender) scanFile(fileName string) ([]url.Values, error) {
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

func (b *newStyleBeaconSender) makeUrlValues(beaconDataMap map[string]string) url.Values {
	result := url.Values{}
	for k, v := range beaconDataMap {
		result.Set(k, v)
	}
	return result
}

func (b *newStyleBeaconSender) httpPost(params url.Values) error {
	headers, err := b.parseHeaders(params)
	if err != nil {
		return fmt.Errorf("parse headers error: %w", err)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%v:%v/beacon/catcher", b.host, b.port), strings.NewReader(params.Encode()))

	if userAgent, ok := headers["User-Agent"]; ok {
		req.Header.Add("User-Agent", userAgent[0])
	} else {
		req.Header.Add("User-Agent", "")
	}

	req.Header.Add("Cf-Ipcountry", b.makeCountryCode(headers))
	req.Header.Add("Cf-Ipcity", b.makeCityName(headers))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.client.Do(req)

	if err != nil {
		return fmt.Errorf("http client error: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("body read error: %w", err)
	}

	log.Println(string(body))

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("expected status: %d, received: %d", http.StatusNoContent, resp.StatusCode)
	}
	return nil
}

func (b *newStyleBeaconSender) parseHeaders(params url.Values) (map[string][]string, error) {
	var headers map[string][]string

	err := json.Unmarshal([]byte(params.Get("request_headers")), &headers)
	if err != nil {
		return nil, fmt.Errorf("unable to parse json[%s]: %w", params, err)
	}

	return headers, nil
}

func (b *newStyleBeaconSender) makeCountryCode(headers map[string][]string) string {
	countryCode := headers["Cf-Ipcountry"][0]
	return countryCode
}

func (b *newStyleBeaconSender) makeCityName(headers map[string][]string) string {
	countryCode := headers["Cf-Ipcity"][0]
	return countryCode
}
