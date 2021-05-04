package main

import (
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/http/server"
)

func main() {
	// Parse the HTTPServer config from environment variables prefixed with HTTP_SERVER (e.g. HTTP_SERVER_PORT)
	envVars := envconfig.PrefixLookuper("HTTP_SERVER_", envconfig.OsLookuper())

	httpServer, err := server.InitializeHTTPServer(envVars)
	if err != nil {
		log.Fatal().Err(err).Msg("HTTPServer initialization failure")
	}
	log.Info().Msg("HTTPServer initialized")

	httpServer.Start()
}
