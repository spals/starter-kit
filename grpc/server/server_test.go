package server_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/grpc/client"
	"github.com/spals/starter-kit/grpc/proto"
	"github.com/spals/starter-kit/grpc/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	grpcPort = 18081
	// The number of milliseconds between checks for the server start
	// NOTE: Increase this number if debugging the server start sequence
	serverStartTickMs = 10
	// The number of milliseconds to wait for the server to start
	// NOTE: Increase this number if debugging the server start sequence
	serverStartTimeoutMs = 50
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
	configMap := make(map[string]string)
	configMap["PORT"] = fmt.Sprintf("%d", grpcPort)
	testLookuper := envconfig.MapLookuper(configMap)

	grpcServer, _ := server.InitializeGrpcServer(testLookuper)
	go func() {
		grpcServer.Start()
	}()

	s.grpcServer = grpcServer

	grpcTarget := fmt.Sprintf("localhost:%d", grpcPort)
	s.grpcClient = client.NewGrpcClient(grpcTarget)
}

func (s *GrpcServerTestSuite) SetupTest() {
	assert := assert.New(s.T())
	// Wait 50 milliseconds for the GrpcServer to be ready
	assert.Eventually(func() bool {
		client := proto.NewHealthClient(s.grpcClient.Conn())
		resp, err := client.GetReady(context.Background(), &proto.ReadyRequest{})
		return err == nil && resp.IsReady
	}, serverStartTimeoutMs*time.Millisecond /*waitFor*/, serverStartTickMs*time.Millisecond /*tick*/)
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
	resp, err := client.GetConfig(context.Background(), &proto.ConfigRequest{})
	if assert.NoError(err) {
		assert.Equal(grpcPort, int(resp.GetConfig().Port))
	}
}
