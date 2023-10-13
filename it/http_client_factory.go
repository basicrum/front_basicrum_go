package it

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func NewHttpClient() *http.Client {
	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:       100,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	return &http.Client{Transport: tr,
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
