package impl

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc/health"
	healthproto "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	ROOT_SERVICE = ""
)

// HealthRegistry ...
// Implementation of auto-generated ConfigServer Grpc framework
// See https://github.com/grpc/grpc/blob/master/doc/health-checking.md
type HealthRegistry struct {
	healthproto.UnimplementedHealthServer

	delegate *health.Server
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
	log.Printf("Marking service '%s' as NOT SERVING", serviceName)
	r.delegate.SetServingStatus(serviceName, healthproto.HealthCheckResponse_NOT_SERVING)
}

// MarkAsServing ...
func (r *HealthRegistry) MarkAsServing(service interface{}) {
	serviceName := serviceName(service)
	log.Printf("Marking service '%s' as SERVING", serviceName)
	r.delegate.SetServingStatus(serviceName, healthproto.HealthCheckResponse_SERVING)
}

// Resume ...
// See health.Server.Resume
func (r *HealthRegistry) Resume() {
	r.delegate.Resume()
}

// Shutdown ...
// See health.Server.Shutdown
func (r *HealthRegistry) Shutdown() {
	r.delegate.Shutdown()
}

// ========== Public Helpers ==========

func serviceName(service interface{}) string {
	return strings.ReplaceAll(fmt.Sprintf("%T", service), "*", "")
}
