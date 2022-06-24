// https://marcofranssen.nl/build-a-go-webserver-on-http-2-using-letsencrypt

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/basicrum/catcher_go/beacon"

	"github.com/basicrum/catcher_go/persistence"

	"github.com/gorilla/mux"
	"github.com/ua-parser/uap-go/uaparser"
)

var (
	domain string
)

const DBNAME = "default"
const DBADDR = "localhost:9000"
const DBAUSERNAME = "default"
const DBPASSWORD = ""
const TABLENAME = "webperf_rum_events"

func main() {
	flag.StringVar(&domain, "domain", "", "domain name to request your certificate")
	flag.Parse()

	// We need to ge the Regexes from here: https://github.com/ua-parser/uap-core/blob/master/regexes.yaml
	uaP, err := uaparser.New("./assets/uaparser_regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Setup the db
	ctx := context.Background()

	err, chConn := persistence.ConnectClickHouse(DBADDR, DBNAME, DBAUSERNAME, DBPASSWORD)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/beacon/catcher", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		b := beacon.FromRequestParams(&r.Form, r.UserAgent(), r.Header)

		re := beacon.ConvertToRumEvent(b, uaP)

		jsonValue, _ := json.Marshal(re)

		persistence.SaveInClickHouse(ctx, chConn, TABLENAME, string(jsonValue))

		if err != nil {
			fmt.Fprint(w, err)
		}
	}).Methods("POST", "GET")

	// fmt.Println("TLS domain", domain)
	// certManager := autocert.Manager{
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist(domain),
	// 	Cache:      autocert.DirCache("certs"),
	// }

	// tlsConfig := certManager.TLSConfig()
	// tlsConfig.GetCertificate = getSelfSignedOrLetsEncryptCert(&certManager)

	// server := http.Server{
	// 	Addr:    initAddress,
	// 	Handler: r,
	// 	// TLSConfig: tlsConfig,
	// }

	// go http.ListenAndServe(":80", certManager.HTTPHandler(nil))

	errdd := http.ListenAndServe(":8087", r)

	if errdd != nil {
		fmt.Println(err)
	}

	// fmt.Println("Server listening on", server.Addr)
	// if err := server.ListenAndServeTLS("", ""); err != nil {
	// 	fmt.Println(err)
	// }
}

// func getSelfSignedOrLetsEncryptCert(certManager *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
// 	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
// 		dirCache, ok := certManager.Cache.(autocert.DirCache)
// 		if !ok {
// 			dirCache = "certs"
// 		}

// 		keyFile := filepath.Join(string(dirCache), hello.ServerName+".key")
// 		crtFile := filepath.Join(string(dirCache), hello.ServerName+".crt")
// 		certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
// 		if err != nil {
// 			fmt.Printf("%s\nFalling back to Letsencrypt\n", err)
// 			return certManager.GetCertificate(hello)
// 		}
// 		fmt.Println("Loaded selfsigned certificate.")
// 		return &certificate, err
// 	}
// }
