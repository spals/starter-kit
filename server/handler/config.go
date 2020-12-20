package handler

import (
	"net/http"

	"starter-kit/server/config"
)

// HTTPServerConfigHandler ...
// HTTP handler which returns HTTP server configuration
type HTTPServerConfigHandler struct {
	config *config.HTTPServerConfig
}

// NewHTTPServerConfigHandler ...
func NewHTTPServerConfigHandler(config *config.HTTPServerConfig) *HTTPServerConfigHandler {
	h := &HTTPServerConfigHandler{config}
	return h
}

func (h *HTTPServerConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(h.config.ToJSONString(false /*prettyPrint*/)))
}
