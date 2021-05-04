package server_test

import (
	"fmt"
	nativelog "log"
	"net/http"
	"testing"
	"time"

	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/http/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	// The number of milliseconds between checks for the server start
	// NOTE: Increase this number if debugging the server start sequence
	serverStartTickMs = 10
	// The number of milliseconds to wait for the server to start
	// NOTE: Increase this number if debugging the server start sequence
	serverStartTimeoutMs = 100
)

// ========== Suite Definition ==========

type HTTPServerTestSuite struct {
	// Extends the testify suite package
	// See https://github.com/stretchr/testify#suite-package
	suite.Suite
	// A reference to the HTTPServer created for testing
	httpServer *server.HTTPServer
	// The base URL to be used by an HTTP client during testing
	httpURLBase string
}

// ========== Setup and Teardown ==========

func (s *HTTPServerTestSuite) SetupSuite() {
	configMap := make(map[string]string)
	testLookuper := envconfig.MapLookuper(configMap)

	httpServer, _ := server.InitializeHTTPServer(testLookuper)
	go func() {
		httpServer.Start()
	}()

	s.httpServer = httpServer
}

func (s *HTTPServerTestSuite) SetupTest() {
	assert := assert.New(s.T())
	// Wait 100 milliseconds for the HTTPServer to be ready
	assert.Eventually(func() bool {
		if s.httpServer.ActivePort() == 0 {
			nativelog.Print("No active port available for HTTP testing")
			return false
		} else if len(s.httpURLBase) == 0 {
			s.httpURLBase = fmt.Sprintf("http://localhost:%d", s.httpServer.ActivePort())
			nativelog.Printf("Using base URL %s for HTTP testing", s.httpURLBase)
		}

		resp, err := http.Get(fmt.Sprintf("%s/ready", s.httpURLBase))
		return err == nil && resp.StatusCode == 200
	}, serverStartTimeoutMs*time.Millisecond /*waitFor*/, serverStartTickMs*time.Millisecond /*tick*/)
}

func (s *HTTPServerTestSuite) TearDownSuite() {
	s.httpServer.Shutdown()
}

// ========== Test Trigger ==========
func TestHTTPServerTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPServerTestSuite))
}

// ========== Tests ==========
func (s *HTTPServerTestSuite) TestGetConfig() {
	assert := assert.New(s.T())
	resp, err := http.Get(fmt.Sprintf("%s/config", s.httpURLBase))
	if assert.NoError(err) {
		assert.Equal(200, resp.StatusCode)
		assert.Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	}
}
