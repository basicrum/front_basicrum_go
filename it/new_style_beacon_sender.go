package it

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
		items, err := scanFile(file)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
	}
	return result, nil
}

func (b *newStyleBeaconSender) httpPost(params url.Values) error {
	var headers map[string][]string

	err := json.Unmarshal([]byte(params.Get("request.headers")), &headers)
	if err != nil {
		return fmt.Errorf("unable to parse json[%s]: %w", params, err)
	}

	countryCode := headers["Cf-Ipcountry"][0]
	cityName := headers["Cf-Ipcity"][0]
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%v:%v/beacon/catcher", b.host, b.port), strings.NewReader(params.Encode()))

	if userAgent, ok := headers["User-Agent"]; ok {
		req.Header.Add("User-Agent", userAgent[0])
	} else {
		req.Header.Add("User-Agent", "")
	}

	req.Header.Add("Cf-Ipcountry", countryCode)
	req.Header.Add("Cf-Ipcity", cityName)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.client.Do(req)

	if err != nil {
		fmt.Println("Client err")
		fmt.Printf("%s", err)
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
