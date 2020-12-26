package main

import (
	"log"

	"github.com/spals/starter-kit/grpc/proto"
	"github.com/spals/starter-kit/grpc/server"
	"github.com/spals/starter-kit/grpc/server/impl"
)

func main() {

	log.Print("Initializing GrpcServer")
	config := &proto.GrpcServerConfig{}
	configServer := impl.NewConfigServer(config)

	grpcServer := server.NewGrpcServer(config, configServer)
	log.Print("GrpcServer initialized")

	grpcServer.Start()
}
