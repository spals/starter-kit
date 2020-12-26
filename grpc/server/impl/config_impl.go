package impl

import (
	"context"

	"github.com/spals/starter-kit/grpc/proto"
)

// ConfigServer ...
type ConfigServer struct {
	proto.UnimplementedConfigServer

	config *proto.GrpcServerConfig
}

// NewConfigServer ...
func NewConfigServer(config *proto.GrpcServerConfig) *ConfigServer {
	s := &ConfigServer{config: config}
	return s
}

// GetConfig ...
func (s *ConfigServer) GetConfig(ctx context.Context, req *proto.ConfigRequest) (*proto.ConfigResponse, error) {
	return &proto.ConfigResponse{Config: s.config}, nil
}
