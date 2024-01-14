package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/basicrum/front_basicrum_go/backup"
	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/service"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

const (
	defaultHTTPPort  = "80"
	defaultHTTPSPort = "443"
	cacheDirPath     = "cache-dir"
)

// Factory is server factory
type Factory struct {
	processService *service.Service
	backupService  backup.IBackup
}

// NewFactory returns server factory
func NewFactory(
	processService *service.Service,
	backupService backup.IBackup,
) *Factory {
	return &Factory{
		processService: processService,
		backupService:  backupService,
	}
}

// Build creates http/https server(s) based on startup configuration
func (f *Factory) Build(sConf config.StartupConfig) ([]*Server, error) {
	httpPort := defaultValue(sConf.Server.Port, defaultHTTPPort)
	httpsPort := defaultValue(sConf.Server.Port, defaultHTTPSPort)

	if !sConf.Server.SSL {
		log.Println("HTTP configuration enabled")
		httpServer := New(
			f.processService,
			f.backupService,
			WithHTTP(httpPort),
		)
		return []*Server{httpServer}, nil
	}

	log.Printf("SSL configuration enabled type[%v]\n", sConf.Server.SSLType)
	switch sConf.Server.SSLType {
	case config.SSLTypeLetsEncrypt:
		allowedHost := sConf.Server.SSLLetsEncrypt.Domain
		log.Printf("SSL Let's Encrypt allowedHost[%v]\n", allowedHost)
		tlsConfig := makeLetsEncryptTLSConfig(allowedHost)
		httpsServer := New(
			f.processService,
			f.backupService,
			WithTLSConfig(defaultHTTPSPort, tlsConfig),
		)
		httpServer := New(
			f.processService,
			f.backupService,
			WithHTTP(httpPort),
		)
		return []*Server{httpsServer, httpServer}, nil
	case config.SSLTypeFile:
		log.Println("SSL files configuration enabled")
		httpsServer := New(
			f.processService,
			f.backupService,
			WithSSL(httpsPort, sConf.Server.SSLFile.SSLFileCertFile, sConf.Server.SSLFile.SSLFileKeyFile),
		)
		httpServer := New(
			f.processService,
			f.backupService,
			WithHTTP(httpPort),
		)
		return []*Server{httpsServer, httpServer}, nil
	default:
		return nil, fmt.Errorf("unsupported SSL type[%v]", sConf.Server.SSLType)
	}
}

func makeLetsEncryptTLSConfig(allowedHost string) *tls.Config {
	client := makeACMEClient()
	m := &autocert.Manager{
		Cache:      autocert.DirCache(cacheDirPath),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(allowedHost),
		Client:     client,
	}
	// nolint: gosec
	return &tls.Config{GetCertificate: m.GetCertificate}
}

func makeACMEClient() *acme.Client {
	directoryURL := os.Getenv("TEST_DIRECTORY_URL")
	if directoryURL == "" {
		return nil
	}
	insecureSkipVerify := os.Getenv("TEST_INSECURE_SKIP_VERIFYy") == "true"
	return &acme.Client{
		DirectoryURL: directoryURL,
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// nolint: gosec
					InsecureSkipVerify: insecureSkipVerify,
				},
			},
		},
	}
}

func defaultValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
