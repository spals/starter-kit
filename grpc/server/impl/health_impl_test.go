package impl_test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/spals/starter-kit/grpc/server/impl"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	healthproto "google.golang.org/grpc/health/grpc_health_v1"
)

// Alias for Protobuf auto-generated Enum type
type ServingStatus = healthproto.HealthCheckResponse_ServingStatus

type mockHealth_WatchServer struct {
	grpc.ServerStream

	ctx           context.Context
	healthUpdates []ServingStatus
	t             *testing.T
}

func newMockHealth_WatchServer(t *testing.T) (mockHealth_WatchServer, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	watchServer := mockHealth_WatchServer{ctx: ctx, t: t}
	return watchServer, cancel
}

func (_m *mockHealth_WatchServer) assertUpdates(expectedUpdates ...ServingStatus) {
	assert := assert.New(_m.t)
	assert.Equal(expectedUpdates, _m.healthUpdates)
}

func (_m *mockHealth_WatchServer) Context() context.Context {
	return _m.ctx
}

func (_m *mockHealth_WatchServer) Send(resp *healthproto.HealthCheckResponse) error {
	log.Printf("WatchServer received status : %s", resp.GetStatus())
	_m.healthUpdates = append(_m.healthUpdates, resp.GetStatus())
	return nil
}

func (_m *mockHealth_WatchServer) startWatch(serviceName string, registry *impl.HealthRegistry, wg *sync.WaitGroup) {
	// Setup a watch in a separate goroutine
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		registry.Watch(&healthproto.HealthCheckRequest{Service: serviceName}, _m) // nolint:errcheck
	}(wg)

	time.Sleep(10 * time.Millisecond) // Wait for the watch to get started
}

func TestBasicCheckServing(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	registry.MarkAsServing(t)

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{Service: "testing.T"})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_SERVING, resp.GetStatus())
	}
}

func TestBasicCheckNotServing(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	registry.MarkAsNotServing(t)

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{Service: "testing.T"})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_NOT_SERVING, resp.GetStatus())
	}
}

func TestBasicWatchServing(t *testing.T) {
	registry := impl.NewHealthRegistry()
	watchServer, cancel := newMockHealth_WatchServer(t)
	var wg sync.WaitGroup
	wg.Add(1)

	watchServer.startWatch("testing.T", registry, &wg)

	registry.MarkAsServing(t)
	time.Sleep(10 * time.Millisecond) // Wait for the watch to fully process
	cancel()

	wg.Wait()
	watchServer.assertUpdates(
		healthproto.HealthCheckResponse_SERVICE_UNKNOWN, // Initial watch value
		healthproto.HealthCheckResponse_SERVING,
	)
}

func TestBasicWatchNotServing(t *testing.T) {
	registry := impl.NewHealthRegistry()
	watchServer, cancel := newMockHealth_WatchServer(t)
	var wg sync.WaitGroup
	wg.Add(1)

	watchServer.startWatch("testing.T", registry, &wg)

	registry.MarkAsNotServing(t)
	time.Sleep(10 * time.Millisecond) // Wait for the watch to fully process
	cancel()

	wg.Wait()
	watchServer.assertUpdates(
		healthproto.HealthCheckResponse_SERVICE_UNKNOWN, // Initial watch value
		healthproto.HealthCheckResponse_NOT_SERVING,
	)
}
