package handler

import (
	"net/http"

	"starter-kit/server/config"
)

// ConfigHandler ...
// HTTP handler which returns HTTP server configuration
type configHandler struct {
	config *config.HTTPServerConfig
}

// NewConfigHandler ...
func NewConfigHandler(config *config.HTTPServerConfig) http.Handler {
	h := &configHandler{config}
	return h
}

func (h *configHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(h.config.ToJSONString()))
}
