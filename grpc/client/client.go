package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	defaultRequestTimeout = 1 * time.Second
)

// GrpcClient ...
type GrpcClient struct {
	ctx    context.Context
	cancel context.CancelFunc

	conn *grpc.ClientConn
}

// NewGrpcClient ...
func NewGrpcClient(target string) *GrpcClient {
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("GrpcClient could not connect: %+v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)

	grpcClient := GrpcClient{conn: conn, ctx: ctx, cancel: cancel}
	return &grpcClient
}

// Close ...
func (c *GrpcClient) Close() error {
	defer c.cancel()
	return c.conn.Close()
}

// Conn ..
func (c *GrpcClient) Conn() *grpc.ClientConn {
	return c.conn
}

// Ctx ...
func (c *GrpcClient) Ctx() context.Context {
	return c.ctx
}
