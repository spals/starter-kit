package main

import (
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/grpc/server"
)

func main() {
	// Parse the GrpcServer config from environment variables prefixed with GRPC_SERVER (e.g. GRPC_SERVER_PORT)
	envVars := envconfig.PrefixLookuper("GRPC_SERVER_", envconfig.OsLookuper())

	grpcServer, err := server.InitializeGrpcServer(envVars)
	if err != nil {
		log.Fatal().Err(err).Msg("GrpcServer initialization failure")
	}
	log.Info().Msg("GrpcServer initialized")

	grpcServer.Start()
}
