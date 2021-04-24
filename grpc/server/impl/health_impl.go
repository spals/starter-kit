package impl

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"google.golang.org/grpc/health"
	healthproto "google.golang.org/grpc/health/grpc_health_v1"
)

// HealthRegistry ...
// Implementation of auto-generated ConfigServer Grpc framework
// See https://github.com/grpc/grpc/blob/master/doc/health-checking.md
type HealthRegistry struct {
	healthproto.UnimplementedHealthServer

	delegate          *health.Server
	mu                sync.RWMutex
	serviceServingMap sync.Map
}

// ========== Constructor ==========

// NewHealthRegistry ...
func NewHealthRegistry() *HealthRegistry {
	healthServer := health.NewServer()

	r := &HealthRegistry{delegate: healthServer}
	return r
}

// ========== Implementation Methods ==========
// These are required implementations based on the grpc service
// definition in grpc_health_v1/health.proto

// Check ...
func (r *HealthRegistry) Check(ctx context.Context, req *healthproto.HealthCheckRequest) (*healthproto.HealthCheckResponse, error) {
	return r.delegate.Check(ctx, req)
}

// Watch ...
func (r *HealthRegistry) Watch(req *healthproto.HealthCheckRequest, stream healthproto.Health_WatchServer) error {
	return r.delegate.Watch(req, stream)
}

// ========== Public Helpers ==========

// MarkAsNotServing ...
func (r *HealthRegistry) MarkAsNotServing(service interface{}) {
	serviceName := serviceName(service)
	r.markServingStatus(serviceName, healthproto.HealthCheckResponse_NOT_SERVING)
}

// MarkAsServing ...
func (r *HealthRegistry) MarkAsServing(service interface{}) {
	serviceName := serviceName(service)
	r.markServingStatus(serviceName, healthproto.HealthCheckResponse_SERVING)
}

// Resume ...
// See health.Server.Resume
func (r *HealthRegistry) Resume() {
	// TODO: Restart rootServiceWatcher
	r.delegate.Resume()
}

// Shutdown ...
// See health.Server.Shutdown
func (r *HealthRegistry) Shutdown() {
	r.delegate.Shutdown()
}

// ========== Private Helpers ==========

func (r *HealthRegistry) markServingStatus(serviceName string, status healthproto.HealthCheckResponse_ServingStatus) {
	serving := (status == healthproto.HealthCheckResponse_SERVING)
	r.serviceServingMap.Store(serviceName, serving)

	r.mu.Lock()
	defer r.mu.Unlock()
	log.Printf("Marking service '%s' as %s", serviceName, status)
	r.delegate.SetServingStatus(serviceName, status)
	r.delegate.SetServingStatus("", r.rootHealthStatus())
}

func (r *HealthRegistry) rootHealthStatus() healthproto.HealthCheckResponse_ServingStatus {
	rootServing := true
	r.serviceServingMap.Range(func(serviceName, serviceServing interface{}) bool {
		rootServing = rootServing && (serviceServing == true)
		return true
	})

	if rootServing {
		return healthproto.HealthCheckResponse_SERVING
	}

	return healthproto.HealthCheckResponse_NOT_SERVING
}

func serviceName(service interface{}) string {
	return strings.ReplaceAll(fmt.Sprintf("%T", service), "*", "")
}
