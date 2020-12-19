package main

import (
	"starter-kit/server"
	"starter-kit/server/config"

	"github.com/sethvargo/go-envconfig"
)

func main() {
	// Parse the HTTPServer config from environment variables prefixed with HTTP_SERVER (e.g. HTTP_SERVER_PORT)
	envVars := envconfig.PrefixLookuper("HTTP_SERVER_", envconfig.OsLookuper())
	c := config.NewHTTPServerConfig(envVars)

	s := server.NewHTTPServer(c)
	s.Start()
}
