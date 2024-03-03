package it

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/stretchr/testify/suite"
)

// Inspired by https://www.gojek.io/blog/golang-integration-testing-made-easy
type e2eTestSuite struct {
	suite.Suite
	dao          *IntegrationDao
	httpSender   *HttpSender
	beaconSender *BeaconSender
}

func (s *e2eTestSuite) SetupTest() {
	var err error
	sConf, err := config.GetStartupConfig()
	s.Assert().NoError(err)
	host := getHost()
	client := NewHttpClient()
	s.httpSender = newHttpSender(
		client,
		host,
		sConf.Server.Port,
	)
	s.beaconSender = newBeaconSender(
		s.httpSender,
	)

	conn, err := dao.NewConnection(
		dao.Server(sConf.Database.Host, sConf.Database.Port, sConf.Database.DatabaseName),
		dao.Auth(sConf.Database.Username, sConf.Database.Password),
	)
	s.Assert().NoError(err)
	s.dao = NewIntegrationDao(
		conn,
		Opts(sConf.Database.TablePrefix),
	)

	s.dao.RecycleTables()
	time.Sleep(10 * time.Second)
}

func getHost() string {
	host := os.Getenv("BRUM_SERVER_HOST")
	if host == "" {
		host = "localhost"
	}
	return host
}

func TestE2ETestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	if os.Getenv("SKIP_E2E") == "true" {
		t.Skip("skipping e2e test")
	}
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) Test_EndToEnd_CountRecords() {
	s.beaconSender.Send("./data/beacon/*.json.lines")
	time.Sleep(2 * time.Second)

	var cntExpect int = 11
	s.Assert().Exactly(cntExpect, s.dao.CountRecords(""))
}

func (s *e2eTestSuite) Test_EndToEnd_BeaconFieldsPersisted() {
	s.beaconSender.Send("./data/misc/all-fields.json.lines")
	time.Sleep(2 * time.Second)

	var cntExpect int = 1

	s.Assert().Exactly(cntExpect, s.dao.CountRecords(""))

	// Set expectations
	cntExpect = 1

	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where cumulative_layout_shift > 0.095"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where first_input_delay = 1"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where user_agent = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where device_type = 'desktop'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where browser_version = '104.0.5112'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where browser_name = 'Chrome'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where ua_vnd = 'Google Inc.'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where ua_plt = 'Win32'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where operating_system = 'Windows'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where operating_system_version = '10'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where device_manufacturer is NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where mob_etype = '4g'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where mob_dl = 10"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where mob_rtt = 50"))
}

func (s *e2eTestSuite) Test_EndToEnd_BeaconFieldsEmpty() {
	s.beaconSender.Send("./data/misc/empty-fields.json.lines")
	time.Sleep(2 * time.Second)

	var cntExpect int = 1

	s.Assert().Exactly(cntExpect, s.dao.CountRecords(""))

	// Set expectations
	cntExpect = 1

	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where cumulative_layout_shift IS NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where first_input_delay IS NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where user_agent IS NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where device_type = 'unknown'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where browser_version IS NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where browser_name = 'Other'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where ua_vnd IS NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where ua_plt IS NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where operating_system = 'Other'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where operating_system_version is NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where device_manufacturer is NULL"))
}

func (s *e2eTestSuite) Test_EndToEnd_BeaconFieldsMissing() {
	s.beaconSender.Send("./data/misc/missing-fields.json.lines")
	time.Sleep(2 * time.Second)

	cntExpect := 1

	s.Assert().Exactly(cntExpect, s.dao.CountRecords(""))

	// Set expectations
	cntExpect = 1

	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where cumulative_layout_shift IS NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where first_input_delay IS NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where user_agent IS NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where device_type = 'unknown'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where browser_version IS NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where browser_name = 'Other'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where ua_vnd IS NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where ua_plt IS NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where operating_system = 'Other'"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where operating_system_version is NULL"))
	s.Assert().Exactly(cntExpect, s.dao.CountRecords("where device_manufacturer is NULL"))
}

func (s *e2eTestSuite) Test_EndToEnd_MobDlFloat() {
	s.beaconSender.Send("./data/misc/mob-dl-float.json.lines")
	time.Sleep(2 * time.Second)

	cntExpect := 2

	s.Assert().Exactly(cntExpect, s.dao.CountRecords(""))

	s.Assert().Exactly(1, s.dao.CountRecords("where mob_dl = 10"))
	s.Assert().Exactly(1, s.dao.CountRecords("where mob_dl = 11"))
}

func (s *e2eTestSuite) Test_EndToEnd_HealthCheck() {
	req, err := http.NewRequest("GET", s.httpSender.BuildUrl("/health"), strings.NewReader(""))
	s.Assert().NoError(err)
	s.httpSender.Send(req, http.StatusOK, "ok")
}
