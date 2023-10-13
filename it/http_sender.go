package it

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HttpSender struct {
	client *http.Client
	host   string
	port   string
}

func newHttpSender(
	client *http.Client,
	host string,
	port string,
) *HttpSender {
	return &HttpSender{
		client: client,
		host:   host,
		port:   port,
	}
}

func (s *HttpSender) Send(req *http.Request, expectedStatusCode int, expectedBody string) {
	err := s.doSend(req, expectedStatusCode, expectedBody)
	if err != nil {
		panic(err)
	}
}

func (s *HttpSender) doSend(req *http.Request, expectedStatusCode int, expectedBody string) error {
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != expectedStatusCode {
		return fmt.Errorf("expected status: %d, received: %d, body: %s", expectedStatusCode, resp.StatusCode, string(body))
	}

	if string(body) != expectedBody {
		return fmt.Errorf("expected body: %s, received: %s", expectedBody, string(body))
	}
	return nil
}

func (s *HttpSender) BuildUrl(path string) string {
	if path != "" {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
	}
	return fmt.Sprintf("http://%v:%v%v", s.host, s.port, path)
}
