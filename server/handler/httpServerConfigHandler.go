package handler

import (
	"net/http"

	"starter-kit/server/config"
)

// HTTP handler which returns HTTP server configuration
type httpServerConfigHandler struct {
	config *config.HTTPServerConfig
}

// NewHTTPServerConfigHandler ...
func NewHTTPServerConfigHandler(config *config.HTTPServerConfig) http.Handler {
	h := &httpServerConfigHandler{config}
	return h
}

func (h *httpServerConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(h.config.ToJSONString()))
}
