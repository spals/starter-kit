package server_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/http/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ========== Suite Definition ==========

type HTTPServerTestSuite struct {
	// Extends the testify suite package
	// See https://github.com/stretchr/testify#suite-package
	suite.Suite
	// The HTTP port that will be used during testing
	httpPort int
	// A reference to the HTTPServer created for testing
	httpServer *server.HTTPServer
	// The base URL to be used by an HTTP client during testing
	httpURLBase string
}

// ========== Setup and Teardown ==========

func (s *HTTPServerTestSuite) SetupSuite() {
	s.httpPort = 18080
	s.httpURLBase = fmt.Sprintf("http://localhost:%d", s.httpPort)

	configMap := make(map[string]string)
	configMap["PORT"] = fmt.Sprintf("%d", s.httpPort)

	testConfig := envconfig.MapLookuper(configMap)
	httpServer, _ := server.InitializeHTTPServer(testConfig)
	go func() {
		httpServer.Start()
	}()

	s.httpServer = httpServer
}

func (s *HTTPServerTestSuite) SetupTest() {
	assert := assert.New(s.T())
	// Wait 50 milliseconds for the HTTPServer to be ready
	assert.Eventually(func() bool {
		resp, err := http.Get(fmt.Sprintf("%s/ready", s.httpURLBase))
		return err == nil && resp.StatusCode == 200
	}, 50*time.Millisecond /*waitFor*/, 10*time.Millisecond /*tick*/)
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
