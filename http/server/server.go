package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
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
			log.Debug().Str("path", "/config").Array("methods", zerolog.Arr().Str("GET")).Msg("Adding HTTP handler")
			router.Path("/config").Methods("GET").Handler(httpServerConfigHandler)

			log.Debug().Str("path", "/live").Array("methods", zerolog.Arr().Str("GET")).Msg("Adding HTTP handler")
			router.Path("/live").Methods("GET").HandlerFunc((*healthCheckHandler).LiveEndpoint)

			log.Debug().Str("path", "/ready").Array("methods", zerolog.Arr().Str("GET")).Msg("Adding HTTP handler")
			router.Path("/ready").Methods("GET").HandlerFunc((*healthCheckHandler).ReadyEndpoint)
		},
	)
	delegate := &http.Server{Handler: router}

	httpServer := &HTTPServer{config, delegate}
	return httpServer
}

// ActivePort ...
// Returns the port on which the server is actively listening.
// This is useful as the server is capable or using a randomly assigned port.
func (s *HTTPServer) ActivePort() int {
	// Note that the port will be re-written in the configuration if a random one is used.
	return s.config.Port
}

// Start ...
func (s *HTTPServer) Start() {
	// Include a graceful server shutdown sequence
	// See https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a#16fd
	httpServerStopped := make(chan os.Signal, 1)
	signal.Notify(httpServerStopped, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	customListener := s.makeCustomListener()

	go func() {
		// If we do not have a custom listener, then use the default listener
		if customListener == nil {
			s.delegate.Addr = fmt.Sprint(":", s.config.Port)
			log.Info().Msgf("HTTPServer listening on port %s", s.delegate.Addr)
			if err := s.delegate.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("HTTPServer start failure with default listener")
			}
		} else {
			log.Info().Msgf("HTTPServer listening on port :%d", customListener.Addr().(*net.TCPAddr).Port)
			if err := s.delegate.Serve(customListener); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("HTTPServer start failure with custom listener")
			}
		}
	}()
	log.Info().Msg("HTTPServer started")

	<-httpServerStopped
	log.Info().Msg("HTTPServer stopped")
	s.Shutdown()
}

// Shutdown ...
func (s *HTTPServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer func() {
		cancel()
	}()

	log.Info().Msg("Shutting down HTTPServer")
	if err := s.delegate.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("HTTPServer shutdown failed")
	}
	log.Info().Msg("HTTPServer shutdown")
}

// If a random port is requested, then make a custom listener on an open port
// Otherwise, return nil
// See https://stackoverflow.com/questions/43424787/how-to-use-next-available-port-in-http-listenandserve
func (s *HTTPServer) makeCustomListener() net.Listener {
	if s.config.Port == 0 {
		log.Debug().Msg("Finding available random port")
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatal().Err(err).Msg("Error while finding random port")
		}

		newPort := listener.Addr().(*net.TCPAddr).Port
		log.Info().Msgf("Overwriting configured port (%d) with random port (%d)", s.config.Port, newPort)
		s.config.Port = newPort
		return listener
	}

	return nil
}

// ========== Private Helpers ==========

// See https://gist.github.com/husobee/fd23681261a39699ee37
type middleware func(http.Handler) http.Handler

func addLoggingMiddleware(
	config *config.HTTPServerConfig,
	rootHandler http.Handler,
) http.Handler {
	return buildChain(
		rootHandler,
		hlog.NewHandler(config.ReqLogger),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("Finished HTTP request")
		}),
		hlog.RemoteAddrHandler("ip"),
		hlog.UserAgentHandler("user_agent"),
		hlog.RefererHandler("referer"),
		hlog.RequestIDHandler("req_id", "Request-Id"),
	)
}

func buildChain(f http.Handler, m ...middleware) http.Handler {
	// If our chain is done, use the original handler
	if len(m) == 0 {
		return f
	}
	// Otherwise nest the handlers
	return m[0](buildChain(f, m[1:cap(m)]...))
}

func makeRouter(
	config *config.HTTPServerConfig,
	routerRegistration routerRegistration,
) http.Handler {
	router := mux.NewRouter()
	routerRegistration(router)

	// Wrap router in a logging handler in order to create access logs
	routerWithLogging := addLoggingMiddleware(config, router)
	return routerWithLogging
}
