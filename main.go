// https://marcofranssen.nl/build-a-go-webserver-on-http-2-using-letsencrypt

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/basicrum/catcher_go/beacon"

	"github.com/gorilla/mux"
	"github.com/ua-parser/uap-go/uaparser"
	"golang.org/x/crypto/acme/autocert"
)

var (
	domain string
)

var initAddress = ":4443"

func getSelfSignedOrLetsEncryptCert(certManager *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		dirCache, ok := certManager.Cache.(autocert.DirCache)
		if !ok {
			dirCache = "certs"
		}

		keyFile := filepath.Join(string(dirCache), hello.ServerName+".key")
		crtFile := filepath.Join(string(dirCache), hello.ServerName+".crt")
		certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			fmt.Printf("%s\nFalling back to Letsencrypt\n", err)
			return certManager.GetCertificate(hello)
		}
		fmt.Println("Loaded selfsigned certificate.")
		return &certificate, err
	}
}

func main() {
	flag.StringVar(&domain, "domain", "", "domain name to request your certificate")
	flag.Parse()

	// We need to ge the Regexes from here: https://github.com/ua-parser/uap-core/blob/master/regexes.yaml
	uaP, err := uaparser.New("./assets/uaparser_regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/beacon/catcher", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		b := beacon.FromRequestParams(&r.Form)
		re := beacon.ConvertToRumEvent(b, uaP)

		jsonValue, _ := json.Marshal(re)

		resp, err := http.Post("http://benthos_app:4195/post", "application/json", bytes.NewBuffer(jsonValue))

		// resp, err := http.PostForm(, urlParams)

		if err != nil {
			fmt.Fprint(w, err)
		}

		fmt.Println(b)

		urlParams := r.URL.Query()

		fmt.Fprint(w, resp)

		fmt.Fprint(w, urlParams)
	}).Methods("POST", "GET")

	fmt.Println("TLS domain", domain)
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache("certs"),
	}

	tlsConfig := certManager.TLSConfig()
	tlsConfig.GetCertificate = getSelfSignedOrLetsEncryptCert(&certManager)
	server := http.Server{
		Addr:      initAddress,
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	fmt.Println("Server listening on", server.Addr)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		fmt.Println(err)
	}
}
