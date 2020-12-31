package impl

import (
	"context"
	"net/http"

	"github.com/heptiolabs/healthcheck"
	"github.com/spals/starter-kit/grpc/proto"
)

// HealthServer ...
// Implementation of auto-generated HealthServer Grpc framework
type HealthServer struct {
	proto.UnimplementedHealthServer

	healthCheckHandler healthcheck.Handler
}

type healthResponseWriter struct {
	http.ResponseWriter

	headers http.Header
	body    []byte
	status  int
}

// ========== Constructor ==========

// NewHealthServer ...
func NewHealthServer(config *proto.GrpcServerConfig) *HealthServer {
	healthCheckHandler := healthcheck.NewHandler()
	configureLivenessChecks(config, healthCheckHandler)
	configureReadinessChecks(config, healthCheckHandler)

	s := &HealthServer{healthCheckHandler: healthCheckHandler}
	return s
}

// ========== Implementation Methods ==========
// These are required implementations based on the grpc service
// definition in config.proto

// GetLive ...
func (s *HealthServer) GetLive(ctx context.Context, req *proto.LiveRequest) (*proto.LiveResponse, error) {
	w := healthResponseWriter{}
	r, _ := http.NewRequest("GET", "local/?full=1", nil)
	s.healthCheckHandler.LiveEndpoint(w, r)

	return nil, nil
}

// GetReady ...
func (s *HealthServer) GetReady(ctx context.Context, req *proto.ReadyRequest) (*proto.ReadyResponse, error) {
	return nil, nil
}

// ========== Private Helpers ==========

func configureLivenessChecks(config *proto.GrpcServerConfig, healthCheckHandler healthcheck.Handler) {
	if config.GetLivenessConfig().GetMaxGoRoutines() > 0 {
		healthCheckHandler.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(int(config.GetLivenessConfig().GetMaxGoRoutines())))
	}
}

func configureReadinessChecks(config *proto.GrpcServerConfig, healthCheckHandler healthcheck.Handler) {

}
