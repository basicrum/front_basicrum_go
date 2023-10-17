// https://marcofranssen.nl/build-a-go-webserver-on-http-2-using-letsencrypt

package main

import (
	"context"
	_ "embed"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/basicrum/front_basicrum_go/backup"
	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/basicrum/front_basicrum_go/geoip"
	"github.com/basicrum/front_basicrum_go/geoip/cloudflare"
	"github.com/basicrum/front_basicrum_go/geoip/maxmind"
	"github.com/basicrum/front_basicrum_go/server"
	"github.com/basicrum/front_basicrum_go/service"
	"github.com/basicrum/front_basicrum_go/service/subscription/caching"
	"github.com/basicrum/front_basicrum_go/service/subscription/disabled"
	"github.com/ua-parser/uap-go/uaparser"
	"golang.org/x/sync/errgroup"
)

//go:embed assets/uaparser_regexes.yaml
var userAgentRegularExpressions []byte

// nolint: revive
func main() {
	sConf, err := config.GetStartupConfig()
	if err != nil {
		log.Fatal(err)
	}

	// We need to get the Regexes from here: https://github.com/ua-parser/uap-core/blob/master/regexes.yaml
	userAgentParser, err := uaparser.NewFromBytes(userAgentRegularExpressions)
	if err != nil {
		log.Fatal(err)
	}

	daoServer := dao.Server(sConf.Database.Host, sConf.Database.Port, sConf.Database.DatabaseName)
	daoAuth := dao.Auth(sConf.Database.Username, sConf.Database.Password)

	conn, err := dao.NewConnection(
		daoServer,
		daoAuth,
	)
	if err != nil {
		log.Fatal(err)
	}

	daoService := dao.New(
		conn,
		dao.Opts(sConf.Database.TablePrefix),
	)

	migrateDaoService := dao.NewMigrationDAO(
		daoServer,
		daoAuth,
		dao.Opts(sConf.Database.TablePrefix),
	)

	err = migrateDaoService.Migrate()
	if err != nil {
		log.Fatalf("migrate database ERROR: %+v", err)
	}

	geopIPService := geoip.NewComposite(
		cloudflare.New(),
		maxmind.New(),
	)

	compressionFactory := backup.NewCompressionWriterFactory(sConf.Backup.Enabled, backup.Compression(sConf.Backup.CompressionType), backup.CompressionLevel(sConf.Backup.CompressionLevel))
	backupInterval := time.Duration(sConf.Backup.IntervalSeconds) * time.Second
	backupService, err := backup.New(sConf.Backup.Enabled, backupInterval, sConf.Backup.Directory, sConf.Backup.ExpiredDirectory, sConf.Backup.UnknownDirectory, compressionFactory)
	if err != nil {
		log.Fatal(err)
	}

	subscriptionService := makeSubscriptionService(sConf, daoService)
	err = subscriptionService.Load()
	if err != nil {
		log.Fatal(err)
	}
	processingService := service.New(
		daoService,
		userAgentParser,
		geopIPService,
		subscriptionService,
		backupService,
	)
	serverFactory := server.NewFactory(processingService, backupService)
	servers, err := serverFactory.Build(*sConf)
	if err != nil {
		log.Fatal(err)
	}

	go processingService.Run()
	startServers(servers)
	if err := stopServers(servers, backupService); err != nil {
		log.Fatalf("Shutdown Failed:%+v", err)
	}
	log.Print("Servers exited properly")
}

func makeSubscriptionService(conf *config.StartupConfig, daoService *dao.DAO) service.ISubscriptionService {
	if !conf.Subscription.Enabled {
		return disabled.New()
	}
	return caching.New(daoService)
}

func startServers(servers []*server.Server) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	for index, srv := range servers {
		go func(srv *server.Server, index int) {
			if err := srv.Serve(); err != nil {
				log.Printf("error start server index[%v] err[%v]\n", index, err)
			}
		}(srv, index)
	}
	log.Print("Servers started")

	<-done
}

func stopServers(servers []*server.Server, backupService backup.IBackup) error {
	log.Print("Stopping servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	g, ctx := errgroup.WithContext(ctx)

	for _, srv := range servers {
		serverCopy := srv
		g.Go(func() error {
			if err := serverCopy.Shutdown(ctx); err != nil {
				log.Printf("Server Shutdown Failed:%+v", err)
				return err
			}
			return nil
		})
	}

	g.Go(func() error {
		backupService.Flush()
		return nil
	})

	// wait for all parallel jobs to finish
	return g.Wait()
}
