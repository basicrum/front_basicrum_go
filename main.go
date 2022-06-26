// https://marcofranssen.nl/build-a-go-webserver-on-http-2-using-letsencrypt

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/rs/cors"

	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/persistence"

	"github.com/ua-parser/uap-go/uaparser"
)

var (
	domain string
)

const TABLENAME = "webperf_rum_events_dummy"

func main() {

	sConf := config.GetStartupConfig()

	flag.StringVar(&domain, "domain", "", "domain name to request your certificate")
	flag.Parse()

	// We need to ge the Regexes from here: https://github.com/ua-parser/uap-core/blob/master/regexes.yaml
	uaP, err := uaparser.New("./assets/uaparser_regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Start: Setup the db
	ctx := context.Background()

	err, chConn := persistence.ConnectClickHouse(sConf.Database.Host, sConf.Database.Port, sConf.Database.DatabaseName, sConf.Database.Username, sConf.Database.Password)
	if err != nil {
		panic(err)
	}
	// End: Setup the db

	mux := http.NewServeMux()

	mux.HandleFunc("/beacon/catcher", func(w http.ResponseWriter, r *http.Request) {
		// @todo: Check if we need to add more response headers
		// access-control-allow-credentials: true
		// access-control-allow-origin: *
		// cache-control: no-cache, no-store, must-revalidate
		// content-length: 0
		// content-type: text/plain
		// cross-origin-resource-policy: cross-origin
		// date: Sat, 25 Jun 2022 10:40:18 GMT
		// expires: Fri, 01 Jan 1990 00:00:00 GMT
		// pragma: no-cache

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")

		w.WriteHeader(http.StatusNoContent)

		defer func(req *http.Request) {
			req.ParseForm()

			b := beacon.FromRequestParams(&req.Form, r.UserAgent(), req.Header)

			re := beacon.ConvertToRumEvent(b, uaP)

			jsonValue, _ := json.Marshal(re)

			persistence.SaveInClickHouse(ctx, chConn, TABLENAME, string(jsonValue))

			if err != nil {
				fmt.Fprint(w, err)
			}
		}(r)
	})

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

	handler := cors.Default().Handler(mux)
	errdd := http.ListenAndServe(":8087", handler)

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
