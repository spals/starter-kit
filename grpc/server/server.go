package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog/log"
	"github.com/spals/starter-kit/grpc/proto"
	"github.com/spals/starter-kit/grpc/server/impl"
	"github.com/spals/starter-kit/grpc/server/logging"
	"google.golang.org/grpc"
	healthproto "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// GrpcServer ...
type GrpcServer struct {
	config         *proto.GrpcServerConfig
	healthRegistry *impl.HealthRegistry // Keep a reference to the health registry so we can shut it down
	delegate       *grpc.Server
}

// NewGrpcServer ...
func NewGrpcServer(
	config *proto.GrpcServerConfig,
	healthRegistry *impl.HealthRegistry,
	configServer *impl.ConfigServer,
) *GrpcServer {
	delegate := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			logging.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			logging.UnaryServerInterceptor(),
		)),
	)

	// Register any service implementations
	reflection.Register(delegate)
	healthproto.RegisterHealthServer(delegate, healthRegistry)
	proto.RegisterConfigServer(delegate, configServer)

	grpcServer := &GrpcServer{config, healthRegistry, delegate}
	return grpcServer
}

// ActivePort ...
// Returns the port on which the server is actively listening.
// This is useful as the server is capable or using a randomly assigned port.
func (s *GrpcServer) ActivePort() int {
	// Note that the port will be re-written in the configuration if a random one is used.
	return int(s.config.GetPort())
}

// Start ...
func (s *GrpcServer) Start() {
	// Include a graceful server shutdown sequence
	// See https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a#16fd
	grpcServerStopped := make(chan os.Signal, 1)
	signal.Notify(grpcServerStopped, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	listener := s.makeListener()
	log.Info().Msgf("GrpcServer listening on port :%d", s.config.GetPort())

	go func() {
		if err := s.delegate.Serve(listener); err != nil {
			log.Fatal().Err(err).Msg("GrpcServer start failure")
		}
	}()
	log.Info().Msg("GrpcServer started")

	<-grpcServerStopped
	log.Info().Msg("GrpcServer stopped")
	s.Shutdown()
}

// Shutdown ...
func (s *GrpcServer) Shutdown() {
	log.Info().Msg("Shutting down GrpcServer")
	s.healthRegistry.Shutdown()
	s.delegate.GracefulStop()
	log.Info().Msg("GrpcServer shutdown")
}

// ========== Private Helpers ==========

func (s *GrpcServer) makeListener() net.Listener {
	// If a random port is requested, then find an open port
	// See https://stackoverflow.com/questions/43424787/how-to-use-next-available-port-in-http-listenandserve
	if s.config.GetPort() == 0 {
		log.Debug().Msg("Finding available random port")
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatal().Err(err).Msg("Error while finding random port")
		}

		newPort := listener.Addr().(*net.TCPAddr).Port
		log.Info().Msgf("Overwriting configured port (%d) with random port (%d)", s.config.Port, newPort)
		s.config.Port = int32(newPort)
		return listener
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GetPort()))
	if err != nil {
		log.Fatal().Err(err).Msgf("Error while listening on port %d", s.config.GetPort())
	}
	return listener
}
