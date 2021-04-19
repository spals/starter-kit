package impl

import (
	"context"

	"github.com/spals/starter-kit/grpc/proto"
)

// ConfigServer ...
// Implementation of auto-generated ConfigServer Grpc framework
type ConfigServer struct {
	proto.UnimplementedConfigServer

	config *proto.GrpcServerConfig
}

// ========== Constructor ==========

// NewConfigServer ...
func NewConfigServer(
	config *proto.GrpcServerConfig,
	healthRegistry *HealthRegistry,
) *ConfigServer {
	s := &ConfigServer{config: config}
	healthRegistry.MarkAsServing(s)

	return s
}

// ========== Implementation Methods ==========
// These are required implementations based on the grpc service
// definition in config.proto

// GetConfig ...
func (s *ConfigServer) GetConfig(ctx context.Context, req *proto.ConfigRequest) (*proto.ConfigResponse, error) {
	return &proto.ConfigResponse{Config: s.config}, nil
}
