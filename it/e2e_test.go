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

	"github.com/basicrum/front_basicrum_go/config"
	"github.com/basicrum/front_basicrum_go/it"
	"github.com/stretchr/testify/suite"
)

// Inspired by https://www.gojek.io/blog/golang-integration-testing-made-easy
type e2eTestSuite struct {
	suite.Suite
	p     *it.Persistence
	sConf *config.StartupConfig
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

func (s *e2eTestSuite) Test_EndToEnd_CountRecords() {
	it.SendBeacons("./data/old_style/*.json", "./data/new_style/*.json.lines")
	time.Sleep(2 * time.Second)

	var cntExpect uint64 = 25
	s.Assert().Exactly(cntExpect, s.p.CountRecords(""))
}

func (s *e2eTestSuite) Test_EndToEnd_BeaconFieldsPersisted() {
	it.SendBeacons("", "./data/misc/all-fields.json.lines")
	time.Sleep(2 * time.Second)

	var cntExpect uint64 = 1

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

	var cntExpect uint64 = 1

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

	var cntExpect uint64 = 1

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
