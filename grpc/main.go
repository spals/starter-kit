package main

import (
	"log"
	"os"

	"github.com/spals/starter-kit/grpc/proto"
	"github.com/spals/starter-kit/grpc/server"
)

func main() {

	log.Print("Initializing GrpcServer")
	config := &proto.GrpcServerConfig{}

	grpcServer, err := server.InitializeGrpcServer(config)
	if err != nil {
		log.Fatalf("GrpcServer initialization failure: %s", err)
		os.Exit(1)
	}
	log.Print("GrpcServer initialized")

	grpcServer.Start()
}
