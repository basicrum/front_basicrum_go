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
	t            *testing.T
	dao          *DAO
	migrationDAO *MigrationDAO
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

	daoServer := Server(sConf.Database.Host, sConf.Database.Port, sConf.Database.DatabaseName)
	daoAuth := Auth(sConf.Database.Username, sConf.Database.Password)

	conn, err := NewConnection(
		daoServer,
		daoAuth,
	)
	s.NoError(err)

	s.dao = New(
		conn,
		Opts(sConf.Database.TablePrefix),
	)

	s.migrationDAO = NewMigrationDAO(
		daoServer,
		daoAuth,
		Opts(sConf.Database.TablePrefix),
	)

	err = s.migrationDAO.Migrate()
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
	s.Equal(1, s.countRows(baseHostsTableName))

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
	s.Equal(1, s.countRows(baseHostsTableName))
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

	// then
	s.Equal(1, s.countRows(baseOwnerHostsTableName))

	// given
	ownerHostname = types.NewOwnerHostname(
		"test1",
		"hostname1",
		types.NewSubscription(time.Now().Add(time.Hour)),
	)

	// when
	err = s.dao.InsertOwnerHostname(ownerHostname)
	s.NoError(err)

	// then
	s.Equal(1, s.countRows(baseOwnerHostsTableName))

	whereClause := "WHERE hostname='hostname1'"
	s.Equal(
		"test1",
		s.selectColumnString("username", baseOwnerHostsTableName, whereClause),
	)
	s.Equal(
		ownerHostname.Subscription.ID,
		s.selectColumnString("subscription_id", baseOwnerHostsTableName, whereClause),
	)
	s.EqualTime(
		ownerHostname.Subscription.ExpiresAt,
		s.selectColumnTime("subscription_expire_at", baseOwnerHostsTableName, whereClause),
	)
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

	// then
	s.Equal(1, s.countRows(baseOwnerHostsTableName))

	// when
	err = s.dao.DeleteOwnerHostname(
		"hostname1",
		"test1",
	)
	s.NoError(err)

	// then
	s.Equal(0, s.countRows(baseOwnerHostsTableName))
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

func (s *daoTestSuite) countRows(tableName string) int {
	s.optimizeFinal(tableName)

	query := fmt.Sprintf("SELECT count(*) FROM %v%v", s.dao.prefix, tableName)
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

func (s *daoTestSuite) optimizeFinal(tableName string) {
	query := fmt.Sprintf("optimize table %v%v final", s.dao.prefix, tableName)
	err := s.dao.conn.Exec(context.Background(), query)
	s.NoError(err)
}

func (s *daoTestSuite) selectColumnString(columnName, tableName, whereClause string) string {
	var value string
	s.selectColumn(columnName, tableName, whereClause, &value)
	return value
}

func (s *daoTestSuite) selectColumnTime(columnName, tableName, whereClause string) time.Time {
	var value time.Time
	s.selectColumn(columnName, tableName, whereClause, &value)
	return value
}

func (s *daoTestSuite) selectColumn(columnName, tableName, whereClause string, value any) {
	s.optimizeFinal(tableName)

	query := fmt.Sprintf("SELECT %v FROM %v%v %v", columnName, s.dao.prefix, tableName, whereClause)
	rows, err := s.dao.conn.Query(context.Background(), query)
	s.NoError(err)
	defer rows.Close()

	if !rows.Next() {
		return
	}
	err = rows.Scan(value)
	s.NoError(err)
}

// EqualTime assert that two times are the same truncated to seconds
func (s *daoTestSuite) EqualTime(expected, actual time.Time) {
	s.Equal(expected.Truncate(time.Second).UnixMilli(), actual.UnixMilli())
}
