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
		// Health Registry
		impl.NewHealthRegistry,
		// Service Implementations
		impl.NewConfigServer,
		// Server
		NewGrpcServer,
	)
	return &GrpcServer{}, nil
}
