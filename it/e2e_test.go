package it_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/it"
	"github.com/basicrum/front_basicrum_go/persistence"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

// Inspired by https://www.gojek.io/blog/golang-integration-testing-made-easy
type e2eTestSuite struct {
	suite.Suite
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

// func (s *e2eTestSuite) SetupSuite() {
// 	serverReady := make(chan bool)

// 	go m.RealMain()
// 	<-serverReady
// }

// func (s *e2eTestSuite) TearDownSuite() {
// 	p, _ := os.FindProcess(syscall.Getpid())
// 	p.Signal(syscall.SIGINT)
// }

// func (s *e2eTestSuite) SetupTest() {
// 	if err := s.dbMigration.Up(); err != nil && err != migrate.ErrNoChange {
// 		s.Require().NoError(err)
// 	}
// }

// func (s *e2eTestSuite) TearDownTest() {
// 	s.NoError(s.dbMigration.Down())
// }

func (s *e2eTestSuite) Test_EndToEnd_CreateArticle() {
	// Start: Setup the db
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	confPath := path + "/config/startup_config.yaml"

	f, err := os.Open(confPath)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	var sConf config.StartupConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&sConf)

	if err != nil {
		log.Println(err)
	}

	p, err := persistence.New(
		persistence.Server(sConf.Database.Host, sConf.Database.Port, sConf.Database.DatabaseName),
		persistence.Auth(sConf.Database.Username, sConf.Database.Password),
	)
	if err != nil {
		log.Fatalf("ERROR: %+v", err)
	}

	p.RecycleTables()
	// End: Setup the db

	time.Sleep(10 * time.Second)

	it.SendBeacons()
	s.NoError(err)

	time.Sleep(2 * time.Second)

	p.CountRecords()
}
