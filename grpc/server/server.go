package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spals/starter-kit/grpc/proto"
	"github.com/spals/starter-kit/grpc/server/impl"
	"google.golang.org/grpc"
)

// GrpcServer ...
type GrpcServer struct {
	config   *proto.GrpcServerConfig
	delegate *grpc.Server
}

// NewGrpcServer ...
func NewGrpcServer(
	config *proto.GrpcServerConfig,
	configServer *impl.ConfigServer,
) *GrpcServer {
	delegate := grpc.NewServer()
	proto.RegisterConfigServer(delegate, configServer)

	grpcServer := &GrpcServer{config, delegate}
	return grpcServer
}

// Start ...
func (s *GrpcServer) Start() {
	// Include a graceful server shutdown sequence
	// See https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a#16fd
	grpcServerStopped := make(chan os.Signal, 1)
	signal.Notify(grpcServerStopped, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	listener := s.makeListener()
	log.Printf("Starting GrpcServer on port %d", s.config.GetPort())

	go func() {
		if err := s.delegate.Serve(listener); err != nil {
			log.Fatalf("GrpcServer start failure: %s", err)
			os.Exit(2)
		}
	}()
	log.Print("GrpcServer started")

	<-grpcServerStopped
	log.Print("GrpcServer stopped")
	s.Shutdown()
}

// Shutdown ...
func (s *GrpcServer) Shutdown() {
	log.Print("Shutting down GrpcServer")
	s.delegate.GracefulStop()
	log.Print("GrpcServer shutdown")
}

func (s *GrpcServer) makeListener() net.Listener {
	// If a random port is requested, then find an open port
	// See https://stackoverflow.com/questions/43424787/how-to-use-next-available-port-in-http-listenandserve
	if s.config.GetAssignRandomPort() || s.config.GetPort() == 0 {
		log.Print("Finding available random port")
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatalf("Error while finding random port: %s", err)
			os.Exit(2)
		}

		newPort := listener.Addr().(*net.TCPAddr).Port
		log.Printf("Overwriting configured port (%d) with random port (%d)", s.config.Port, newPort)
		s.config.Port = int32(newPort)
		return listener
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GetPort()))
	if err != nil {
		log.Fatalf("Error while listening on port %d: %s", s.config.GetPort(), err)
		os.Exit(2)
	}
	return listener
}
