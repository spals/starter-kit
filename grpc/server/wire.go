// +build wireinject

package server

import (
	"github.com/google/wire"
	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/grpc/server/config"
	"github.com/spals/starter-kit/grpc/server/impl"
)

// InitializeGrpcServer ...
func InitializeGrpcServer(l envconfig.Lookuper) (*GrpcServer, error) {
	wire.Build(
		// Configuration
		config.NewGrpcServerConfig,
		// Service Implementations
		impl.NewConfigServer,
		impl.NewHealthServer,
		// Server
		NewGrpcServer,
	)
	return &GrpcServer{}, nil
}
