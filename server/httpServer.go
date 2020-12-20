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

	"starter-kit/server/config"
	"starter-kit/server/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// HTTPServer ...
type HTTPServer struct {
	config   *config.HTTPServerConfig
	delegate *http.Server
}

// NewHTTPServer ...
func NewHTTPServer(config *config.HTTPServerConfig) *HTTPServer {
	router := makeRouter(config)
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
		log.Printf("Starting HTTPServer with default listener: %s", s.config.ToJSONString())
	} else {
		log.Printf("Starting HTTPServer with custom listener: %s", s.config.ToJSONString())
	}

	go func() {
		// If we do not have a custom listener, then use the default listener
		if customListener == nil {
			s.delegate.Addr = fmt.Sprint(":", s.config.Port)
			if err := s.delegate.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPServer start failure with default listener: %s", err)
				os.Exit(1)
			}
		} else {
			if err := s.delegate.Serve(customListener); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPServer start failure with custom listener: %s", err)
			}
		}
	}()
	log.Print("HTTPServer started")

	<-httpServerStopped
	log.Print("HTTPServer stopped")
	s.stop()
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
			os.Exit(1)
		}

		newPort := listener.Addr().(*net.TCPAddr).Port
		log.Printf("Overwriting configured port (%d) with random port (%d)", s.config.Port, newPort)
		s.config.Port = newPort
		return listener
	}

	return nil
}

func (s *HTTPServer) stop() {
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

// ========== Private Helpers ==========

func makeRouter(config *config.HTTPServerConfig) http.Handler {
	router := mux.NewRouter()
	registerRoutes(config, router)

	// Wrap router in a logging handler in order to create access logs
	// See https://godoc.org/github.com/gorilla/handlers#LoggingHandler
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	return loggedRouter
}

func registerRoutes(config *config.HTTPServerConfig, router *mux.Router) {
	configHandler := handler.NewConfigHandler(config)
	router.Path("/config").Methods("GET").Handler(configHandler)

	healthCheckHandler := handler.NewHealthCheckHandler(config)
	router.Path("/live").Methods("GET").HandlerFunc(healthCheckHandler.LiveEndpoint)
	router.Path("/ready").Methods("GET").HandlerFunc(healthCheckHandler.ReadyEndpoint)
}
