package server

import (
	"context"
	"fmt"
	"log"
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
	addr := fmt.Sprint(":", config.Port)
	router := makeRouter(config)

	delegate := &http.Server{Addr: addr, Handler: router}

	httpServer := &HTTPServer{config, delegate}
	return httpServer
}

// Start ...
func (s *HTTPServer) Start() {
	log.Printf("Starting HTTPServer with %s", s.config.ToJSONString())

	// Include a graceful server shutdown sequence
	// See https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a#16fd
	httpServerStopped := make(chan os.Signal, 1)
	signal.Notify(httpServerStopped, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.delegate.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTPServer start failure: %s", err)
		}
	}()
	log.Printf("HTTPServer started at %+v", s.delegate.Addr)

	<-httpServerStopped
	log.Print("HTTPServer stopped")
	s.stop()
}

func (s *HTTPServer) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer func() {
		cancel()
	}()

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

	defaultHandler := handler.NewDefaultHandler()
	router.Path("/").Handler(defaultHandler)
}
