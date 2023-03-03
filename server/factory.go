package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/basicrum/front_basicrum_go/backup"
	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/service"
	"golang.org/x/crypto/acme/autocert"
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
	if !sConf.Server.SSL {
		httpServer := New(
			sConf.Server.Port,
			f.processService,
			f.backupService,
		)
		return []*Server{httpServer}, nil
	}

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
		tlsConfig := &tls.Config{
			GetCertificate: m.GetCertificate,
			MinVersion:     tls.VersionTLS12,
		}
		httpsServer := New(
			sConf.Server.Port,
			f.processService,
			f.backupService,
			WithTLSConfig(tlsConfig),
		)
		httpServer := New(
			sConf.Server.SSLLetsEncrypt.Port,
			f.processService,
			f.backupService,
			WithHandlerAdapter(m.HTTPHandler),
		)
		return []*Server{httpsServer, httpServer}, nil
	case config.SSLTypeFile:
		httpsServer := New(
			sConf.Server.Port,
			f.processService,
			f.backupService,
			WithSSLFile(sConf.Server.SSLFile.SSLFileCertFile, sConf.Server.SSLFile.SSLFileCertFile),
		)
		return []*Server{httpsServer}, nil
	default:
		return nil, fmt.Errorf("unsupported ssl type[%v]", sConf.Server.SSLType)
	}
}
