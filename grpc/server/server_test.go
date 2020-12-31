package server_test

import (
	"fmt"
	"testing"

	"github.com/spals/starter-kit/grpc/client"
	"github.com/spals/starter-kit/grpc/proto"
	"github.com/spals/starter-kit/grpc/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	grpcPort = 18081
)

var (
	grpcTarget = fmt.Sprintf("localhost:%d", grpcPort)
)

// ========== Suite Definition ==========

type GrpcServerTestSuite struct {
	// Extends the testify suite package
	// See https://github.com/stretchr/testify#suite-package
	suite.Suite
	// A reference to the GrpcServer created for testing
	grpcServer *server.GrpcServer
	// A reference to the GrpcClient created for testing
	grpcClient *client.GrpcClient
}

// ========== Setup and Teardown ==========

func (s *GrpcServerTestSuite) SetupSuite() {
	testConfig := proto.GrpcServerConfig{Port: grpcPort}
	grpcServer, _ := server.InitializeGrpcServer(&testConfig)
	go func() {
		grpcServer.Start()
	}()

	s.grpcServer = grpcServer
	s.grpcClient = client.NewGrpcClient(grpcTarget)
}

func (s *GrpcServerTestSuite) SetupTest() {
	// assert := assert.New(s.T())
	// // Wait 50 milliseconds for the HTTPServer to be ready
	// assert.Eventually(func() bool {
	// 	resp, err := http.Get(fmt.Sprintf("%s/ready", s.httpURLBase))
	// 	return err == nil && resp.StatusCode == 200
	// }, 50*time.Millisecond /*waitFor*/, 10*time.Millisecond /*tick*/)
}

func (s *GrpcServerTestSuite) TearDownSuite() {
	s.grpcClient.Close()
	s.grpcServer.Shutdown()
}

// ========== Test Trigger ==========
func TestGrpcServerTestSuite(t *testing.T) {
	suite.Run(t, new(GrpcServerTestSuite))
}

// ========== Tests ==========
func (s *GrpcServerTestSuite) TestGetConfig() {
	assert := assert.New(s.T())

	client := proto.NewConfigClient(s.grpcClient.Conn())
	resp, err := client.GetConfig(s.grpcClient.Ctx(), &proto.ConfigRequest{})
	if assert.NoError(err) {
		assert.False(resp.GetConfig().AssignRandomPort)
		assert.Equal(grpcPort, resp.GetConfig().Port)
	}
}
