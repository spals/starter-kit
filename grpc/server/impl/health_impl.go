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
	w := newHealthResponseWriter()
	r, _ := http.NewRequest("GET", buildHealthURL(req.GetFull()), nil)
	s.healthCheckHandler.LiveEndpoint(w, r)

	resp := proto.LiveResponse{IsLive: w.status == http.StatusOK}
	return &resp, nil
}

// GetReady ...
func (s *HealthServer) GetReady(ctx context.Context, req *proto.ReadyRequest) (*proto.ReadyResponse, error) {
	w := newHealthResponseWriter()
	r, _ := http.NewRequest("GET", buildHealthURL(req.GetFull()), nil)
	s.healthCheckHandler.ReadyEndpoint(w, r)

	resp := proto.ReadyResponse{IsReady: w.status == http.StatusOK}
	return &resp, nil
}

// ========== Private Helpers ==========

func buildHealthURL(isFull bool) string {
	url := "local/?full="
	if isFull {
		url += "1"
	} else {
		url += "0"
	}

	return url
}

func configureLivenessChecks(config *proto.GrpcServerConfig, healthCheckHandler healthcheck.Handler) {
	// if config.GetLivenessConfig().GetMaxGoRoutines() > 0 {
	// 	healthCheckHandler.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(int(config.GetLivenessConfig().GetMaxGoRoutines())))
	// }
}

func configureReadinessChecks(config *proto.GrpcServerConfig, healthCheckHandler healthcheck.Handler) {

}

type healthResponseWriter struct {
	http.ResponseWriter

	headers http.Header
	body    []byte
	status  int
}

func newHealthResponseWriter() *healthResponseWriter {
	return &healthResponseWriter{headers: make(http.Header)}
}

func (w *healthResponseWriter) Header() http.Header {
	return w.headers
}

func (w *healthResponseWriter) Write(body []byte) (int, error) {
	w.body = body
	return len(body), nil
}

func (w *healthResponseWriter) WriteHeader(status int) {
	w.status = status
}
