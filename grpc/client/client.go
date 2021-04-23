package client

import (
	"log"

	"google.golang.org/grpc"
)

// NewGrpcClientConnForTest ...
// Creates a Grpc client connection that can be used for testing purposes
func NewGrpcClientConnForTest(target string) *grpc.ClientConn {
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("GrpcClient could not connect: %+v", err)
	}

	return conn
}
