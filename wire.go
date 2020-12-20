// +build wireinject

package main

import (
	"starter-kit/server"
	"starter-kit/server/config"
	"starter-kit/server/handler"

	"github.com/google/wire"
	"github.com/sethvargo/go-envconfig"
)

// InitializeHTTPServer ...
func InitializeHTTPServer(l envconfig.Lookuper) (*server.HTTPServer, error) {
	wire.Build(
		// Configuration
		config.NewHTTPServerConfig,
		// Handlers
		handler.NewHealthCheckHandler,
		handler.NewHTTPServerConfigHandler,
		// Server
		server.NewHTTPServer,
	)
	return &server.HTTPServer{}, nil
}
