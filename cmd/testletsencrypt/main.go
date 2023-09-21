package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	prepareDomainIP()
	client := makeHTTPClient()
	urlAddress := fmt.Sprintf("https://%s/health", domain())

	found := false
	for i := 1; i <= 3; i++ {
		resp, err := client.Get(urlAddress)

		// then
		if err != nil {
			log.Printf("try[%v] received error[%v]\n", i, err)
			continue
		}
		log.Printf("try[%v] received status code[%v]\n", i, resp.StatusCode)

		_ = resp.Body.Close()
		if http.StatusOK == resp.StatusCode {
			found = true
			log.Printf("try[%v] found\n", i)
			break
		}
	}

	if !found {
		os.Exit(1)
	}
}

func makeHTTPClient() http.Client {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				// nolint: gosec
				InsecureSkipVerify: true,
			},
		},
	}
	return client
}

func domain() string {
	return osGetenvRequired("DOMAIN")
}

func domainIP() string {
	return osGetenvRequired("DOMAIN_IP")
}

func prepareDomainIP() {
	requireNoError(appendToFile("/etc/hosts", fmt.Sprintf("%s %s\n", domainIP(), domain())))
}

func appendToFile(filename string, lines string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	go func() {
		if err := f.Close(); err != nil {
			log.Print(err)
		}
	}()

	_, err = f.WriteString(lines)
	return err
}

func osGetenvRequired(name string) string {
	value := os.Getenv(name)
	if value == "" {
		// nolint: revive
		log.Fatalf("required[%v]\n", name)
	}
	return value
}

func requireNoError(err error) {
	if err != nil {
		// nolint: revive
		log.Fatal(err)
	}
}
