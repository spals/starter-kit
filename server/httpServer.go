package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"starter-kit/server/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// HTTPServer ...
type HTTPServer struct {
	delegate *http.Server
}

// NewHTTPServer ...
func NewHTTPServer(port int16) *HTTPServer {
	addr := fmt.Sprint(":", port)
	router := makeRouter()

	delegate := &http.Server{Addr: addr, Handler: router}

	httpServer := &HTTPServer{delegate}
	return httpServer
}

// Start ...
func (s *HTTPServer) Start() {
	log.Printf("Starting HTTPServer...")

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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer func() {
		cancel()
	}()

	if err := s.delegate.Shutdown(ctx); err != nil {
		log.Fatalf("HTTPServer shutdown failed: %+v", err)
	}
	log.Print("HTTPServer shutdown")
}

// ========== Private Helpers ==========

func makeRouter() http.Handler {
	router := mux.NewRouter()
	registerRoutes(router)

	// Wrap router in a logging handler in order to create access logs
	// See https://godoc.org/github.com/gorilla/handlers#LoggingHandler
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	return loggedRouter
}

func registerRoutes(router *mux.Router) {
	defaultHandler := handler.NewDefaultHandler()
	router.Path("/").Handler(defaultHandler)
}
