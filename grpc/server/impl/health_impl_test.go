package impl_test

import (
	"context"
	"testing"

	"github.com/spals/starter-kit/grpc/server/impl"
	"github.com/stretchr/testify/assert"
	healthproto "google.golang.org/grpc/health/grpc_health_v1"
)

// const (
// 	grpcPort = 54321
// )

// ========== Suite Definition ==========

// type HealthRegistryTestSuite struct {
// 	// Extends the testify suite package
// 	// See https://github.com/stretchr/testify#suite-package
// 	suite.Suite
// 	// A reference to the GrpcServer created for testing
// 	grpcServer     *grpc.Server
// 	healthRegistry *impl.HealthRegistry
// }

// Fake service definitions used for testing
// type TestService struct{}

// ========== Setup and Teardown ==========

// func (s *HealthRegistryTestSuite) SetupSuite() {
// 	assert := assert.New(s.T())
// 	grpcServer := grpc.NewServer()
// 	healthRegistry := impl.NewHealthRegistry()
// 	healthproto.RegisterHealthServer(grpcServer, healthRegistry)

// 	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
// 	assert.NoErrorf(err, "Failed to create listener on port %d", grpcPort)

// 	go func() {
// 		grpcServer.Serve(listener)
// 	}()

// 	s.grpcServer = grpcServer
// 	s.healthRegistry = healthRegistry
// }

// func (s *HealthRegistryTestSuite) TearDownSuite() {
// 	s.healthRegistry.Shutdown()
// 	s.grpcServer.Stop()
// }

// ========== Test Trigger ==========
// func TestHealthRegistryTestSuite(t *testing.T) {
// 	suite.Run(t, new(HealthRegistryTestSuite))
// }

// ========== Tests ==========
// func (s *HealthRegistryTestSuite) TestBasicCheckServing() {
// 	assert := assert.New(s.T())
//
//  registry := impl.NewHealthRegistry()
//  registry.MarkAsServing(s.T())
func TestBasicCheckServing(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	registry.MarkAsServing(t)

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{Service: "testing.T"})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_SERVING, resp.GetStatus())
	}
}

// func (s *HealthRegistryTestSuite) TestBasicCheckNotServing() {
// 	assert := assert.New(s.T())

// 	registry := impl.NewHealthRegistry()
// 	registry.MarkAsNotServing(s.T())

func TestBasicCheckNotServing(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	registry.MarkAsNotServing(t)

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{Service: "testing.T"})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_NOT_SERVING, resp.GetStatus())
	}
}

// func (s *HealthRegistryTestSuite) TestWatch() {
// 	assert := assert.New(s.T())

// 	testService := TestService{}

// 	conn := s.createGrpcConn()
// 	client := healthproto.NewHealthClient(conn)
// 	req := &healthproto.HealthCheckRequest{Service: "impl_test.TestService"}
// 	stream, err := client.Watch(context.Background(), req)
// 	assert.NoError(err)

// 	done := make(chan bool)

// 	go func() {
// 		for {
// 			log.Printf("Receiving on stream")
// 			resp, err := stream.Recv()
// 			if err == io.EOF {
// 				done <- true //means stream is finished
// 				return
// 			}
// 			if err != nil {
// 				log.Fatalf("cannot receive %v", err)
// 			}
// 			log.Printf("Resp received: %s", resp.GetStatus())
// 		}
// 	}()
// 	log.Printf("Calling markAsServing")
// 	s.healthRegistry.MarkAsServing(testService)
// 	log.Printf("Calling markAsServing")
// 	s.healthRegistry.MarkAsServing(testService)
// 	log.Printf("Calling markAsNotServing")
// 	s.healthRegistry.MarkAsNotServing(testService)
// 	log.Printf("Calling markAsNotServing")
// 	s.healthRegistry.MarkAsNotServing(testService)
// 	conn.Close()

// 	<-done
// 	log.Printf("Finished")
// }

// ========== Private Helpers ==========

// func (s *HealthRegistryTestSuite) createGrpcConn() *grpc.ClientConn {
// 	assert := assert.New(s.T())
// 	conn, err := grpc.Dial(fmt.Sprintf(":%d", grpcPort), grpc.WithInsecure())
// 	assert.NoErrorf(err, "Could not connect to Grpc server on port %d", grpcPort)

// 	return conn
// }
