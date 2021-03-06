package server_test

import (
	"context"
	"fmt"
	nativelog "log"
	"sync"
	"testing"
	"time"

	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/grpc/client"
	"github.com/spals/starter-kit/grpc/proto"
	"github.com/spals/starter-kit/grpc/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	healthproto "google.golang.org/grpc/health/grpc_health_v1"
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

type GrpcServerTestSuite struct {
	// Extends the testify suite package
	// See https://github.com/stretchr/testify#suite-package
	suite.Suite
	// A reference to the GrpcServer created for testing
	grpcServer *server.GrpcServer
	// A reference to the Grpc connection created for testing
	grpcConn *grpc.ClientConn
}

// ========== Setup and Teardown ==========

func (s *GrpcServerTestSuite) SetupSuite() {
	configMap := make(map[string]string)
	configMap["LOG_LEVEL"] = "trace"
	testLookuper := envconfig.MapLookuper(configMap)

	grpcServer, _ := server.InitializeGrpcServer(testLookuper)
	go func() {
		grpcServer.Start()
	}()

	s.grpcServer = grpcServer
}

func (s *GrpcServerTestSuite) SetupTest() {
	assert := assert.New(s.T())
	// Wait 100 milliseconds for the GrpcServer to be ready
	assert.Eventually(func() bool {
		if s.grpcServer.ActivePort() == 0 {
			nativelog.Print("No active port available for Grpc testing")
			return false
		} else if s.grpcConn == nil {
			grpcTarget := fmt.Sprintf("localhost:%d", s.grpcServer.ActivePort())
			s.grpcConn = client.NewGrpcClientConnForTest(grpcTarget)
			nativelog.Printf("Using target %s for Grpc testing", grpcTarget)
		}

		healthClient := healthproto.NewHealthClient(s.grpcConn)
		resp, err := healthClient.Check(context.Background(), &healthproto.HealthCheckRequest{})
		return err == nil && resp.GetStatus() == healthproto.HealthCheckResponse_SERVING
	}, serverStartTimeoutMs*time.Millisecond /*waitFor*/, serverStartTickMs*time.Millisecond /*tick*/)
}

func (s *GrpcServerTestSuite) TearDownSuite() {
	s.grpcConn.Close()
	s.grpcServer.Shutdown()
}

// ========== Test Trigger ==========
func TestGrpcServerTestSuite(t *testing.T) {
	suite.Run(t, new(GrpcServerTestSuite))
}

// ========== Tests ==========
func (s *GrpcServerTestSuite) TestGetConfig() {
	assert := assert.New(s.T())

	configClient := proto.NewConfigClient(s.grpcConn)
	resp, err := configClient.GetConfig(context.Background(), &proto.ConfigRequest{})
	if assert.NoError(err) {
		assert.Equal(s.grpcServer.ActivePort(), int(resp.GetConfig().Port))
	}
}

func (s *GrpcServerTestSuite) TestHealthCheck() {
	assert := assert.New(s.T())

	healthClient := healthproto.NewHealthClient(s.grpcConn)
	resp, err := healthClient.Check(context.Background(), &healthproto.HealthCheckRequest{Service: "impl.ConfigServer"})

	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_SERVING, resp.GetStatus())
	}
}

func (s *GrpcServerTestSuite) TestHealthWatch() {
	assert := assert.New(s.T())

	healthClient := healthproto.NewHealthClient(s.grpcConn)
	ctx, cancel := context.WithCancel(context.Background())
	healthWatchClient, err := healthClient.Watch(ctx, &healthproto.HealthCheckRequest{Service: "impl.ConfigServer"})
	assert.NoError(err)

	var healthUpdates []healthproto.HealthCheckResponse_ServingStatus
	var wg sync.WaitGroup
	wg.Add(1)
	go func(healthUpdates *[]healthproto.HealthCheckResponse_ServingStatus, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			resp, err := healthWatchClient.Recv()
			if err != nil { // Error is expected when the context is cancelled
				break
			}
			*healthUpdates = append(*healthUpdates, resp.GetStatus()) // nolint: staticcheck
		}
	}(&healthUpdates, &wg)

	time.Sleep(10 * time.Millisecond) // Wait for the watch to get started
	cancel()                          // Signal to watchServer to stop

	wg.Wait()
	assert.Equal(
		[]healthproto.HealthCheckResponse_ServingStatus{
			healthproto.HealthCheckResponse_SERVING,
		},
		healthUpdates,
	)
}
