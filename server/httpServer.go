package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	log.Printf("Starting HTTPServer at %+v", s.delegate.Addr)
	httpError := s.delegate.ListenAndServe()
	if httpError != nil {
		log.Fatal("FAILURE!")
	}

	// s.stop()
}

// func (s *HTTPServer) stop() {
// 	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

// 	if err := s.delegate.Shutdown(ctx); err != nil {
// 		log.Fatalf("HTTPServer shutdown failed: %+v", err)
// 	}
// 	log.Print("HTTPServer shutdown")
// }

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
