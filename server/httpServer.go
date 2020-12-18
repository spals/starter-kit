package server

import (
	"fmt"
	"log"
	"net/http"

	"starter-kit/server/handler"

	"github.com/gorilla/mux"
)

// HTTPServer ...
type HTTPServer struct {
	delegate *http.Server
}

// NewHTTPServer ...
func NewHTTPServer(port int16) *HTTPServer {
	addr := fmt.Sprint(":", port)
	router := mux.NewRouter()
	registerRoutes(router)

	delegate := &http.Server{Addr: addr, Handler: router}

	httpServer := &HTTPServer{delegate}
	return httpServer
}

func registerRoutes(router *mux.Router) {
	defaultHandler := handler.NewDefaultHandler()
	router.Path("/").Handler(defaultHandler)
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
