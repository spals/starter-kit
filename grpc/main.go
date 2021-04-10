package main

import (
	"log"
	"os"

	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/grpc/server"
)

func main() {
	// Parse the GrpcServer config from environment variables prefixed with GRPC_SERVER (e.g. GRPC_SERVER_PORT)
	envVars := envconfig.PrefixLookuper("GRPC_SERVER_", envconfig.OsLookuper())

	log.Print("Initializing GrpcServer")
	grpcServer, err := server.InitializeGrpcServer(envVars)
	if err != nil {
		log.Fatalf("GrpcServer initialization failure: %s", err)
		os.Exit(1)
	}
	log.Print("GrpcServer initialized")

	grpcServer.Start()
}
