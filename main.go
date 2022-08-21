// https://marcofranssen.nl/build-a-go-webserver-on-http-2-using-letsencrypt

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/basicrum/front_basicrum_go/backup"
	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/persistence"
	"github.com/eapache/go-resiliency/batcher"
	"github.com/rs/cors"
)

var (
	domain string
)

func main() {

	sConf := config.GetStartupConfig()

	flag.StringVar(&domain, "domain", "", "domain name to request your certificate")
	flag.Parse()

	backupInterval := time.Duration(sConf.Backup.IntervalSeconds) * time.Second

	b := batcher.New(backupInterval, func(params []interface{}) error {
		if sConf.Backup.Enabled {
			backup.Do(params, sConf.Backup.Directory)
		}
		return nil
	})

	p, err := persistence.New(
		persistence.Server(sConf.Database.Host, sConf.Database.Port, sConf.Database.DatabaseName),
		persistence.Auth(sConf.Database.Username, sConf.Database.Password),
		persistence.Opts(sConf.Database.TablePrefix),
	)

	if err != nil {
		log.Fatalf("ERROR: %+v", err)
	}

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
			if parseErr := r.ParseForm(); err != nil {
				log.Println(parseErr)
			}

			// We need this in case we would like to re-import beacons
			// Also created_at is used for event date when we persist data in the DB
			if !r.Form.Has("created_at") {
				r.Form.Set("created_at", time.Now().UTC().Format("2006-01-02 15:04:05"))
			}

			// Persist Event in ClickHouse
			p.Save(req, "webperf_rum_events")

			// Archiving logic
			if sConf.Backup.Enabled {
				forArchiving := r.Form

				// Flatten headers later
				h, hErr := json.Marshal(r.Header)

				if hErr != nil {
					log.Println(hErr)
				}

				forArchiving.Add("request_headers", string(h))

				go b.Run(forArchiving)
			}
		}(r)
	})

	// log.Println("TLS domain", domain)
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

	log.Println("Starting the server on port: " + sConf.Server.Port)

	handler := cors.Default().Handler(mux)
	errdd := http.ListenAndServe(":"+sConf.Server.Port, handler)

	if errdd != nil {
		log.Println(errdd)
	}

	// log.Println("Server listening on", server.Addr)
	// if err := server.ListenAndServeTLS("", ""); err != nil {
	// 	log.Println(err)
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
// 			log.Printf("%s\nFalling back to Letsencrypt\n", err)
// 			return certManager.GetCertificate(hello)
// 		}
// 		log.Println("Loaded selfsigned certificate.")
// 		return &certificate, err
// 	}
// }
