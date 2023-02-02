// https://marcofranssen.nl/build-a-go-webserver-on-http-2-using-letsencrypt

// nolint: cyclop
package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"

	"github.com/basicrum/front_basicrum_go/backup"
	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/local/github.com/eapache/go-resiliency/batcher"
	"github.com/basicrum/front_basicrum_go/persistence"
	"github.com/rs/cors"
	"github.com/ua-parser/uap-go/uaparser"
)

//go:embed assets/uaparser_regexes.yaml
var uaRegexes []byte

// nolint: funlen, revive, gocognit
func main() {
	sConf, err := config.GetStartupConfig()
	if err != nil {
		log.Fatal(err)
	}

	// @TODO: Move uaP dependency outside the persistance
	// We need to get the Regexes from here: https://github.com/ua-parser/uap-core/blob/master/regexes.yaml
	uaP, err := uaparser.NewFromBytes(uaRegexes)
	if err != nil {
		log.Fatal(err)
	}

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
		uaP,
	)

	if err != nil {
		log.Fatalf("ERROR: %+v", err)
	}

	err = p.CreateTable()
	if err != nil {
		log.Fatalf("create table ERROR: %+v", err)
	}

	go p.Run()

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

		// Prep for Async work
		parseErr := r.ParseForm()
		if parseErr != nil {
			log.Println(parseErr)
			return
		}

		f := r.Form
		h := r.Header
		uaStr := r.UserAgent()

		// We need this in case we would like to re-import beacons
		// Also created_at is used for event date when we persist data in the DB
		if !f.Has("created_at") {
			f.Set("created_at", time.Now().UTC().Format("2006-01-02 15:04:05"))
		}

		// Persist Event in ClickHouse
		go func() {
			p.Events <- p.Event(&f, &h, uaStr)
		}()

		// Archiving logic
		if sConf.Backup.Enabled {
			forArchiving := f

			// Flatten headers later
			h, hErr := json.Marshal(h)

			if hErr != nil {
				log.Println(hErr)
			}

			forArchiving.Add("request_headers", string(h))

			go func(forArchiving url.Values) {
				if err := b.Run(forArchiving); err != nil {
					log.Printf("Error archiving url[%v] err[%v]", forArchiving, err)
				}
			}(forArchiving)
		}
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte("ok"))
	})
	log.Println("Starting the server on port: " + sConf.Server.Port)

	handler := cors.Default().Handler(mux)

	server := &http.Server{
		Addr:    ":" + sConf.Server.Port,
		Handler: handler,
		// https://deepsource.io/directory/analyzers/go/issues/GO-S2114
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {

		// nolint: nestif
		if sConf.Server.SSL {
			log.Printf("SSL configuration enabled type[%v]\n", sConf.Server.SSLType)
			switch sConf.Server.SSLType {
			case config.SSLTypeLetsEncrypt:
				dataDir := os.TempDir()
				allowedHost := sConf.Server.SSLLetsEncrypt.Domain
				log.Printf("SSL allowedHost[%v]\n", allowedHost)
				hostPolicy := func(ctx context.Context, host string) error {
					if host == allowedHost {
						return nil
					}
					return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
				}
				m := &autocert.Manager{
					Prompt:     autocert.AcceptTOS,
					HostPolicy: hostPolicy,
					Cache:      autocert.DirCache(dataDir),
				}
				server.TLSConfig = &tls.Config{
					GetCertificate: m.GetCertificate,
					MinVersion:     tls.VersionTLS12,
				}

				g, _ := errgroup.WithContext(context.Background())
				g.Go(func() error {
					log.Printf("starting https server on port[%v]", sConf.Server.Port)
					return server.ListenAndServeTLS("", "")
				})
				httpServer := &http.Server{
					Addr:    ":" + sConf.Server.SSLLetsEncrypt.Port,
					Handler: m.HTTPHandler(handler),
					// https://deepsource.io/directory/analyzers/go/issues/GO-S2114
					ReadHeaderTimeout: 3 * time.Second,
					ReadTimeout:       5 * time.Second,
					WriteTimeout:      5 * time.Second,
					IdleTimeout:       120 * time.Second,
				}
				g.Go(func() error {
					log.Printf("starting http server on port[%v]", sConf.Server.SSLLetsEncrypt.Port)
					return httpServer.ListenAndServe()
				})
				if err := g.Wait(); err != nil {
					log.Println(err)
				}
			case config.SSLTypeFile:
				log.Printf("starting https server on port[%v]", sConf.Server.Port)
				errdd := server.ListenAndServe()
				if errdd != nil {
					log.Println(errdd)
				}
			default:
				log.Fatalf("unsupported ssl type[%v]", sConf.Server.SSLType)
			}
		} else {
			log.Printf("starting http server on port[%v]", sConf.Server.Port)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server Shutdown Failed:%+v", err)
			return err
		}
		return nil
	})

	g.Go(func() error {
		b.Flush()
		return nil
	})

	// wait for all parallel jobs to finish
	if err := g.Wait(); err != nil {
		// nolint: gocritic
		log.Fatalf("Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
