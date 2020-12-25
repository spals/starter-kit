package main

import (
	"log"
	"os"

	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/server"
)

func main() {
	// Parse the HTTPServer config from environment variables prefixed with HTTP_SERVER (e.g. HTTP_SERVER_PORT)
	envVars := envconfig.PrefixLookuper("HTTP_SERVER_", envconfig.OsLookuper())

	log.Print("Initializing HTTPServer")
	httpServer, err := server.InitializeHTTPServer(envVars)
	if err != nil {
		log.Fatalf("HTTPServer initialization failure: %s", err)
		os.Exit(1)
	}
	log.Print("HTTPServer initialized")

	httpServer.Start()
}
