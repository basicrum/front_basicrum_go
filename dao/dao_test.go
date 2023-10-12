package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/stretchr/testify/suite"
)

const sleepDuration = 2 * time.Second

type daoTestSuite struct {
	suite.Suite
	t   *testing.T
	dao *DAO
}

func TestDaoTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	testSuite := new(daoTestSuite)
	testSuite.t = t
	suite.Run(t, testSuite)
}

func (s *daoTestSuite) SetupTest() {
	config.SetTestDefaultConfig()
	sConf, err := config.GetStartupConfig()
	s.NoError(err)

	s.dao, err = New(
		Server(sConf.Database.Host, sConf.Database.Port, sConf.Database.DatabaseName),
		Auth(sConf.Database.Username, sConf.Database.Password),
		Opts(sConf.Database.TablePrefix),
	)
	s.NoError(err)

	err = s.dao.Migrate()
	s.NoError(err)

	s.deleteAll()
}

func (s *daoTestSuite) deleteAll() {
	s.truncateTable(baseTableName)
	s.truncateTable(baseHostsTableName)
}

func (s *daoTestSuite) TearDownTest() {
	if s.dao != nil {
		_ = s.dao.Close()
	}
}

func (s *daoTestSuite) Test_SaveHost() {
	// given
	event := beacon.NewHostnameEvent(
		"host1",
		"2022-08-27 05:53:00",
	)

	// when
	err := s.dao.SaveHost(event)
	s.NoError(err)
	// and
	sleep()

	// then
	s.Equal(1, s.countHosts(baseHostsTableName))

	// given
	event = beacon.NewHostnameEvent(
		"host1",
		"2022-08-27 06:53:00",
	)

	// when
	err = s.dao.SaveHost(event)
	s.NoError(err)
	// and
	sleep()

	// then
	s.Equal(1, s.countHosts(baseHostsTableName))
}

func (s *daoTestSuite) Test_InsertOwnerHostname() {
	// given
	ownerHostname := types.NewOwnerHostname(
		"test1",
		"hostname1",
		types.NewSubscription(time.Now()),
	)

	// when
	err := s.dao.InsertOwnerHostname(ownerHostname)
	s.NoError(err)
	// and
	sleep()

	// then
	s.Equal(1, s.countHosts(baseOwnerHostsTableName))

	// given
	ownerHostname = types.NewOwnerHostname(
		"test1",
		"hostname1",
		types.NewSubscription(time.Now().Add(time.Hour)),
	)

	// when
	err = s.dao.InsertOwnerHostname(ownerHostname)
	s.NoError(err)
	// and
	sleep()

	// then
	s.Equal(1, s.countHosts(baseOwnerHostsTableName))
}

func (s *daoTestSuite) Test_DeleteOwnerHostname() {
	// given
	ownerHostname := types.NewOwnerHostname(
		"test1",
		"hostname1",
		types.NewSubscription(time.Now()),
	)

	// when
	err := s.dao.InsertOwnerHostname(ownerHostname)
	s.NoError(err)
	// and
	sleep()

	// then
	s.Equal(1, s.countHosts(baseOwnerHostsTableName))

	// when
	err = s.dao.DeleteOwnerHostname(
		"hostname1",
		"test1",
	)
	s.NoError(err)
	// and
	sleep()

	// then
	s.Equal(0, s.countHosts(baseOwnerHostsTableName))
}

func sleep() {
	time.Sleep(sleepDuration)
}

func (s *daoTestSuite) truncateTable(table string) {
	conn := s.dao.conn
	prefix := s.dao.prefix
	dropQuery := fmt.Sprintf("TRUNCATE TABLE %s%s", prefix, table)
	err := conn.Exec(context.Background(), dropQuery)
	s.NoError(err)
}

func (s *daoTestSuite) countHosts(hostsTable string) int {
	s.optimizeFinal()

	query := fmt.Sprintf("SELECT count(*) FROM %v%v ", s.dao.prefix, hostsTable)
	rows, err := s.dao.conn.Query(context.Background(), query)
	s.NoError(err)
	defer rows.Close()

	if !rows.Next() {
		return 0
	}

	var result uint64

	err = rows.Scan(&result)
	s.NoError(err)

	return int(result)
}

func (s *daoTestSuite) optimizeFinal() {
	query := fmt.Sprintf("optimize table %v%v final", s.dao.prefix, baseHostsTableName)
	err := s.dao.conn.Exec(context.Background(), query)
	s.NoError(err)
}
