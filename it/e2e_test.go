package it_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/basicrum/front_basicrum_go/it"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/suite"
)

// Inspired by https://www.gojek.io/blog/golang-integration-testing-made-easy

type e2eTestSuite struct {
	suite.Suite
	dbConnectionStr string
	port            int
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
	it.SendBeacons()

	reqStr := `{"title":"e2eTitle", "content": "e2eContent", "author":"e2eauthor"}`
	req, err := http.NewRequest(echo.POST, "http://localhost:8087/beacon/catcher", strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusNoContent, response.StatusCode)

	byteBody, err := ioutil.ReadAll(response.Body)
	s.NoError(err)

	s.Equal("", strings.Trim(string(byteBody), "\n"))
	response.Body.Close()
}
