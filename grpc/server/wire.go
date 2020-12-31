// +build wireinject

package server

import (
	"github.com/google/wire"
	"github.com/spals/starter-kit/grpc/proto"
	"github.com/spals/starter-kit/grpc/server/impl"
)

// InitializeGrpcServer ...
func InitializeGrpcServer(config *proto.GrpcServerConfig) (*GrpcServer, error) {
	wire.Build(
		// Service Implementations
		impl.NewConfigServer,
		impl.NewHealthServer,
		// Server
		NewGrpcServer,
	)
	return &GrpcServer{}, nil
}
