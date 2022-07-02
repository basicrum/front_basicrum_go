package it_test

import (
	"context"
	"testing"
	"time"

	"github.com/basicrum/front_basicrum_go/it"
	"github.com/basicrum/front_basicrum_go/persistence"
	"github.com/stretchr/testify/suite"
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
	ctx := context.Background()

	err, chConn := persistence.ConnectClickHouse(
		"localhost",
		"9000",
		"default",
		"default",
		"")
	if err != nil {
		panic(err)
	}

	persistence.RecycleTables(ctx, chConn)
	// End: Setup the db

	it.SendBeacons()
	s.NoError(err)

	time.Sleep(2 * time.Second)

	persistence.CountRecords(ctx, chConn)
}
