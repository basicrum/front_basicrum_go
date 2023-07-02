package it_test

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"
	"time"

	"github.com/basicrum/front_basicrum_go/beacon"
	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/dao"
	"github.com/basicrum/front_basicrum_go/it"
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/ua-parser/uap-go/uaparser"
)

// Inspired by https://www.gojek.io/blog/golang-integration-testing-made-easy
type e2eTestSuite struct {
	suite.Suite
	p     *it.Persistence
	sConf *config.StartupConfig
	dao   *dao.DAO
}

func (s *e2eTestSuite) SetupTest() {
	var err error
	s.sConf, err = config.GetStartupConfig()
	if err != nil {
		log.Fatal(err)
	}
	sConf := s.sConf

	s.p, err = it.New(
		it.Server(sConf.Database.Host, sConf.Database.Port, sConf.Database.DatabaseName, sConf.Database.TablePrefix),
		it.Auth(sConf.Database.Username, sConf.Database.Password),
		it.Opts(sConf.Database.TablePrefix),
	)
	if err != nil {
		log.Fatalf("ERROR: %+v", err)
	}
	s.dao, err = dao.New(
		dao.Server(sConf.Database.Host, sConf.Database.Port, sConf.Database.DatabaseName),
		dao.Auth(sConf.Database.Username, sConf.Database.Password),
		dao.Opts(sConf.Database.TablePrefix),
	)
	if err != nil {
		log.Fatalf("ERROR: %+v", err)
	}

	s.p.RecycleTables()
	// End: Setup the db

	time.Sleep(10 * time.Second)
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

func (s *e2eTestSuite) Test_EndToEnd_CountRecords1() {
	it.SendBeacons("./data/old_style/1638405781.json", "")
	time.Sleep(2 * time.Second)

	assert.Equal(s.T(), 4, s.p.CountRecords(""))
}

type mockGeoIPService struct{}

func (*mockGeoIPService) CountryAndCity(header http.Header, ipString string) (string, string, error) {
	return "", "", nil
}

type mockUserAgentParser struct{}

func (*mockUserAgentParser) Parse(line string) *uaparser.Client {
	return &uaparser.Client{
		UserAgent: &uaparser.UserAgent{},
		Os:        &uaparser.Os{},
		Device:    &uaparser.Device{},
	}
}

func (s *e2eTestSuite) Test_DAO_Save() {
	theBeacon := beacon.Beacon{
		CreatedAt: "2021-12-02 00:42:27",
		Mob_Dl:    "2.7",
		Mob_Rtt:   "5.4",
		Mob_Etype: "5G",
		Rt_Sl:     "0",
	}
	re := beacon.ConvertToRumEvent(theBeacon, &types.Event{}, &mockUserAgentParser{}, &mockGeoIPService{})
	err := s.dao.Save(re)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), 1, s.p.CountRecords(""))
}

func (s *e2eTestSuite) Test_EndToEnd_CountRecords2() {
	it.SendBeacons("./data/old_style/1638405962.json", "")
	time.Sleep(2 * time.Second)

	assert.Equal(s.T(), 4, s.p.CountRecords(""))
}

func (s *e2eTestSuite) Test_EndToEnd_CountRecords3() {
	it.SendBeacons("./data/old_style/1638406081.json", "")
	time.Sleep(2 * time.Second)

	assert.Equal(s.T(), 2, s.p.CountRecords(""))
}

func (s *e2eTestSuite) Test_EndToEnd_CountRecords4() {
	it.SendBeacons("./data/old_style/1638406141.json", "")
	time.Sleep(2 * time.Second)

	assert.Equal(s.T(), 4, s.p.CountRecords(""))
}

func (s *e2eTestSuite) Test_EndToEnd_CountRecords21() {
	it.SendBeacons("", "./data/new_style/ab.json.lines")
	time.Sleep(2 * time.Second)

	assert.Equal(s.T(), 11, s.p.CountRecords(""))
}

func (s *e2eTestSuite) Test_EndToEnd_BeaconFieldsPersisted() {
	it.SendBeacons("", "./data/misc/all-fields.json.lines")
	time.Sleep(2 * time.Second)

	var cntExpect int = 1

	s.Assert().Exactly(cntExpect, s.p.CountRecords(""))

	// Set expectations
	cntExpect = 1

	s.Assert().Exactly(cntExpect, s.p.CountRecords("where cumulative_layout_shift > 0.095"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where first_input_delay = 1"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where user_agent = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where device_type = 'desktop'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where browser_version = '104.0.5112'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where browser_name = 'Chrome'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where ua_vnd = 'Google Inc.'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where ua_plt = 'Win32'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where operating_system = 'Windows'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where operating_system_version = '10'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where device_manufacturer is NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where mob_etype = '4g'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where mob_dl = 10"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where mob_rtt = 50"))
}

func (s *e2eTestSuite) Test_EndToEnd_BeaconFieldsEmpty() {
	it.SendBeacons("", "./data/misc/empty-fields.json.lines")
	time.Sleep(2 * time.Second)

	var cntExpect int = 1

	s.Assert().Exactly(cntExpect, s.p.CountRecords(""))

	// Set expectations
	cntExpect = 1

	s.Assert().Exactly(cntExpect, s.p.CountRecords("where cumulative_layout_shift IS NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where first_input_delay IS NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where user_agent IS NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where device_type = 'unknown'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where browser_version IS NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where browser_name = 'Other'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where ua_vnd IS NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where ua_plt IS NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where operating_system = 'Other'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where operating_system_version is NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where device_manufacturer is NULL"))
}

func (s *e2eTestSuite) Test_EndToEnd_BeaconFieldsMissing() {
	it.SendBeacons("", "./data/misc/missing-fields.json.lines")
	time.Sleep(2 * time.Second)

	var cntExpect int = 1

	s.Assert().Exactly(cntExpect, s.p.CountRecords(""))

	// Set expectations
	cntExpect = 1

	s.Assert().Exactly(cntExpect, s.p.CountRecords("where cumulative_layout_shift IS NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where first_input_delay IS NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where user_agent IS NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where device_type = 'unknown'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where browser_version IS NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where browser_name = 'Other'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where ua_vnd IS NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where ua_plt IS NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where operating_system = 'Other'"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where operating_system_version is NULL"))
	s.Assert().Exactly(cntExpect, s.p.CountRecords("where device_manufacturer is NULL"))
}

func (s *e2eTestSuite) Test_EndToEnd_HealthCheck() {
	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:       100,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr,
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%v:%v/health", s.sConf.Server.Host, s.sConf.Server.Port), strings.NewReader(""))

	resp, err := client.Do(req)

	if err != nil {
		log.Println("Client err")
		log.Printf("%s", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Body read err")
		log.Printf("%s", err)
	}

	s.Assert().Exactly(200, resp.StatusCode)
	s.Assert().Exactly("ok", string(body))
}
