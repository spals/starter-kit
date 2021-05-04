package impl_test

import (
	"context"
	nativelog "log"
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
	nativelog.Printf("MockHealthWatcher received status : %s", resp.GetStatus())
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

func TestCheckBasicServing(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsServing(t)

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{Service: "testing.T"})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_SERVING, resp.GetStatus())
	}
}

func TestCheckBasicNotServing(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsNotServing(t)

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{Service: "testing.T"})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_NOT_SERVING, resp.GetStatus())
	}
}

func TestRootHealthBasicServing(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsServing(t)

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_SERVING, resp.GetStatus())
	}
}

func TestRootHealthBasicNotServing(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsNotServing(t)

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_NOT_SERVING, resp.GetStatus())
	}
}

func TestRootHealthMultiServing(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsServing(t)
	registry.MarkAsServing("")

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_SERVING, resp.GetStatus())
	}
}

func TestRootHealthMultiNotServing(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsNotServing(t) // Any not serving service means the root service is not serving
	registry.MarkAsServing("")

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_NOT_SERVING, resp.GetStatus())
	}
}

func TestRootHealthMultiNotServing2(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsServing(t)
	registry.MarkAsNotServing("") // Any not serving service means the root service is not serving

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_NOT_SERVING, resp.GetStatus())
	}
}

func TestRootHealthTransitionToServing(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsNotServing(t)
	registry.MarkAsServing(t)

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_SERVING, resp.GetStatus())
	}
}

func TestRootHealthTransitionToNotServing(t *testing.T) {
	assert := assert.New(t)

	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsServing(t)
	registry.MarkAsNotServing(t)

	resp, err := registry.Check(context.Background(), &healthproto.HealthCheckRequest{})
	if assert.NoError(err) {
		assert.Equal(healthproto.HealthCheckResponse_NOT_SERVING, resp.GetStatus())
	}
}

func TestWatchBasicServing(t *testing.T) {
	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsServing(t)

	mockHealthWatcher, cancel := newMockHealth_WatchServer(t)
	var wg sync.WaitGroup
	wg.Add(1)

	mockHealthWatcher.startWatch("testing.T", registry, &wg)
	cancel() // Signal to watchServer to stop

	wg.Wait()
	mockHealthWatcher.assertUpdates(
		healthproto.HealthCheckResponse_SERVING,
	)
}

func TestWatchBasicNotServing(t *testing.T) {
	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsNotServing(t)

	mockHealthWatcher, cancel := newMockHealth_WatchServer(t)
	var wg sync.WaitGroup
	wg.Add(1)

	mockHealthWatcher.startWatch("testing.T", registry, &wg)
	cancel() // Signal to watchServer to stop

	wg.Wait()
	mockHealthWatcher.assertUpdates(
		healthproto.HealthCheckResponse_NOT_SERVING,
	)
}

func TestWatchIgnoreDupStatus(t *testing.T) {
	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsNotServing(t)

	mockHealthWatcher, cancel := newMockHealth_WatchServer(t)
	var wg sync.WaitGroup
	wg.Add(1)

	mockHealthWatcher.startWatch("testing.T", registry, &wg)
	registry.MarkAsServing(t) // Update with same serving status twice
	registry.MarkAsServing(t)
	time.Sleep(10 * time.Millisecond) // Wait for the watch to fully process
	cancel()                          // Signal to watchServer to stop

	wg.Wait()
	mockHealthWatcher.assertUpdates(
		healthproto.HealthCheckResponse_NOT_SERVING,
		healthproto.HealthCheckResponse_SERVING, // NOTE: Only one SERVING status received -- dups are ignored
	)
}

func TestWatchIgnorePreWatchStatusChanges(t *testing.T) {
	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsNotServing(t)
	registry.MarkAsServing(t) // Update serving status prior to watch

	mockHealthWatcher, cancel := newMockHealth_WatchServer(t)
	var wg sync.WaitGroup
	wg.Add(1)

	mockHealthWatcher.startWatch("testing.T", registry, &wg)
	cancel() // Signal to watchServer to stop

	wg.Wait()
	mockHealthWatcher.assertUpdates(
		healthproto.HealthCheckResponse_SERVING, // NOTE: Only one status received -- status changes prior to watch are ignored
	)
}

func TestWatchMultiStatus(t *testing.T) {
	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsServing(t)

	mockHealthWatcher, cancel := newMockHealth_WatchServer(t)
	var wg sync.WaitGroup
	wg.Add(1)

	mockHealthWatcher.startWatch("testing.T", registry, &wg)
	registry.MarkAsNotServing(t)      // Update serving status after watch start
	time.Sleep(10 * time.Millisecond) // Wait for the watch to fully process
	cancel()                          // Signal to watchServer to stop

	wg.Wait()
	mockHealthWatcher.assertUpdates(
		healthproto.HealthCheckResponse_SERVING,
		healthproto.HealthCheckResponse_NOT_SERVING,
	)
}

func TestWatchMultiWatch(t *testing.T) {
	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()
	registry.MarkAsServing(t)
	registry.MarkAsServing("")

	mockHealthWatcher, cancel := newMockHealth_WatchServer(t)
	var wg sync.WaitGroup
	wg.Add(2)

	mockHealthWatcher.startWatch("testing.T", registry, &wg)
	mockHealthWatcher.startWatch("string", registry, &wg)
	cancel() // Signal to watchServer to stop

	wg.Wait()
	mockHealthWatcher.assertUpdates(
		healthproto.HealthCheckResponse_SERVING,
		healthproto.HealthCheckResponse_SERVING,
	)
}

func TestWatchUnknownService(t *testing.T) {
	registry := impl.NewHealthRegistry()
	defer registry.Shutdown()

	mockHealthWatcher, cancel := newMockHealth_WatchServer(t)
	var wg sync.WaitGroup
	wg.Add(1)

	mockHealthWatcher.startWatch("testing.T", registry, &wg)
	registry.MarkAsServing(t)         // Register service after watch start
	time.Sleep(10 * time.Millisecond) // Wait for the watch to fully process
	cancel()                          // Signal to watchServer to stop

	wg.Wait()
	mockHealthWatcher.assertUpdates(
		healthproto.HealthCheckResponse_SERVICE_UNKNOWN, // Initial state on watch start
		healthproto.HealthCheckResponse_SERVING,
	)
}
