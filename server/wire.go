// +build wireinject

package server

import (
	"starter-kit/server/config"
	"starter-kit/server/handler"

	"github.com/google/wire"
	"github.com/sethvargo/go-envconfig"
)

// InitializeHTTPServer ...
func InitializeHTTPServer(l envconfig.Lookuper) (*HTTPServer, error) {
	wire.Build(
		// Configuration
		config.NewHTTPServerConfig,
		// Handlers
		handler.NewHealthCheckHandler,
		handler.NewHTTPServerConfigHandler,
		// Server
		NewHTTPServer,
	)
	return &HTTPServer{}, nil
}
