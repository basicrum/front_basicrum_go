package server

import (
	"fmt"
	"log"

	"github.com/basicrum/front_basicrum_go/backup"
	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/service"
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
			sConf.PrivateAPI.Token,
			WithHTTP(httpPort),
		)
		return []*Server{httpServer}, nil
	}

	log.Printf("SSL configuration enabled type[%v]\n", sConf.Server.SSLType)
	switch sConf.Server.SSLType {
	case config.SSLTypeLetsEncrypt:
		allowedHost := sConf.Server.SSLLetsEncrypt.Domain
		log.Printf("SSL Let's Encrypt allowedHost[%v]\n", allowedHost)
		httpsServer := New(
			f.processService,
			f.backupService,
			sConf.PrivateAPI.Token,
			WithListener(autocert.NewListener(allowedHost)),
		)
		httpServer := New(
			f.processService,
			f.backupService,
			sConf.PrivateAPI.Token,
			WithHTTP(httpPort),
		)
		return []*Server{httpsServer, httpServer}, nil
	case config.SSLTypeFile:
		log.Println("SSL files configuration enabled")
		httpsServer := New(
			f.processService,
			f.backupService,
			sConf.PrivateAPI.Token,
			WithSSL(httpsPort, sConf.Server.SSLFile.SSLFileCertFile, sConf.Server.SSLFile.SSLFileKeyFile),
		)
		httpServer := New(
			f.processService,
			f.backupService,
			sConf.PrivateAPI.Token,
			WithHTTP(httpPort),
		)
		return []*Server{httpsServer, httpServer}, nil
	default:
		return nil, fmt.Errorf("unsupported SSL type[%v]", sConf.Server.SSLType)
	}
}

func defaultValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
