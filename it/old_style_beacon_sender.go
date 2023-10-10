package it

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

const countryCode = "DE"

type oldStyleBeaconSender struct {
	client *http.Client
	host   string
	port   string
}

func newOldStyleBeaconSender(
	client *http.Client,
	host string,
	port string,
) *oldStyleBeaconSender {
	return &oldStyleBeaconSender{
		client: client,
		host:   host,
		port:   port,
	}
}

func (b *oldStyleBeaconSender) Send(path string) {
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

func (b *oldStyleBeaconSender) readFiles(path string) ([]url.Values, error) {
	files, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	result := []url.Values{}
	for _, file := range files {
		items, err := processFile(file)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
	}
	return result, nil
}

func (b *oldStyleBeaconSender) httpPost(params url.Values) error {
	log.Println(strings.NewReader(params.Encode()))

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%v:%v/beacon/catcher", b.host, b.port), strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	req.Header.Add("User-Agent", makeUserAgent(params))
	req.Header.Add("CF-IPCountry", countryCode)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("response error: %w", err)
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

func makeUserAgent(params url.Values) string {
	userAgent := params.Get("user.agent")

	if len(userAgent) > 0 {
		return userAgent
	}
	return "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36"
}
