// +build wireinject

package server

import (
	"github.com/google/wire"
	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/server/config"
	"github.com/spals/starter-kit/server/handler"
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
