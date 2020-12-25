package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	"github.com/spals/starter-kit/http/server/config"
	"github.com/spals/starter-kit/http/server/handler"
)

// HTTPServer ...
type HTTPServer struct {
	config   *config.HTTPServerConfig
	delegate *http.Server
}

// Function callback definition used to register routes in HTTPServer router
type routerRegistration func(router *mux.Router)

// NewHTTPServer ...
// Create a new HTTPServer with the given configuration and request handlers.
//
// New request handlers should be added here and in wire.go and registered in
// the router below.
func NewHTTPServer(
	config *config.HTTPServerConfig,
	healthCheckHandler *healthcheck.Handler,
	httpServerConfigHandler *handler.HTTPServerConfigHandler,
) *HTTPServer {
	router := makeRouter(
		config,
		func(router *mux.Router) {
			router.Path("/config").Methods("GET").Handler(httpServerConfigHandler)
			router.Path("/live").Methods("GET").HandlerFunc((*healthCheckHandler).LiveEndpoint)
			router.Path("/ready").Methods("GET").HandlerFunc((*healthCheckHandler).ReadyEndpoint)
		},
	)
	delegate := &http.Server{Handler: router}

	httpServer := &HTTPServer{config, delegate}
	return httpServer
}

// Start ...
func (s *HTTPServer) Start() {
	// Include a graceful server shutdown sequence
	// See https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a#16fd
	httpServerStopped := make(chan os.Signal, 1)
	signal.Notify(httpServerStopped, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	customListener := s.makeCustomListener()
	if customListener == nil {
		log.Printf("Starting HTTPServer on port %d", s.config.Port)
	} else {
		log.Printf("Starting HTTPServer on port %d", customListener.Addr().(*net.TCPAddr).Port)
	}

	go func() {
		// If we do not have a custom listener, then use the default listener
		if customListener == nil {
			s.delegate.Addr = fmt.Sprint(":", s.config.Port)
			if err := s.delegate.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPServer start failure with default listener: %s", err)
				os.Exit(2)
			}
		} else {
			if err := s.delegate.Serve(customListener); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPServer start failure with custom listener: %s", err)
				os.Exit(2)
			}
		}
	}()
	log.Print("HTTPServer started")

	<-httpServerStopped
	log.Print("HTTPServer stopped")
	s.Shutdown()
}

// Shutdown ...
func (s *HTTPServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer func() {
		cancel()
	}()

	log.Print("Shutting down HTTPServer")
	if err := s.delegate.Shutdown(ctx); err != nil {
		log.Fatalf("HTTPServer shutdown failed: %+v", err)
	}
	log.Print("HTTPServer shutdown")
}

// If a random port is requested, then make a custom listener on an open port
// Otherwise, return nil
// See https://stackoverflow.com/questions/43424787/how-to-use-next-available-port-in-http-listenandserve
func (s *HTTPServer) makeCustomListener() net.Listener {
	if s.config.AssignRandomPort {
		log.Print("Finding available random port")
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatalf("Error while finding random port: %s", err)
			os.Exit(2)
		}

		newPort := listener.Addr().(*net.TCPAddr).Port
		log.Printf("Overwriting configured port (%d) with random port (%d)", s.config.Port, newPort)
		s.config.Port = newPort
		return listener
	}

	return nil
}

// ========== Private Helpers ==========

func makeRouter(
	config *config.HTTPServerConfig,
	routerRegistration routerRegistration,
) http.Handler {
	router := mux.NewRouter()
	routerRegistration(router)

	// Wrap router in a logging handler in order to create access logs
	// See https://godoc.org/github.com/gorilla/handlers#LoggingHandler
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	return loggedRouter
}
