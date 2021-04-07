package client

import (
	"log"

	"google.golang.org/grpc"
)

// GrpcClient ...
type GrpcClient struct {
	conn *grpc.ClientConn
}

// NewGrpcClient ...
func NewGrpcClient(target string) *GrpcClient {
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("GrpcClient could not connect: %+v", err)
	}

	grpcClient := GrpcClient{conn: conn}
	return &grpcClient
}

// Close ...
func (c *GrpcClient) Close() error {
	return c.conn.Close()
}

// Conn ..
func (c *GrpcClient) Conn() *grpc.ClientConn {
	return c.conn
}
