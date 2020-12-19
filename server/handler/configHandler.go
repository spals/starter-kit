package handler

import (
	"net/http"

	"starter-kit/server/config"
)

// ConfigHandler ...
// HTTP handler which returns HTTP server configuration
type ConfigHandler struct {
	config *config.HTTPServerConfig
}

// NewConfigHandler ...
func NewConfigHandler(config *config.HTTPServerConfig) *ConfigHandler {
	h := &ConfigHandler{config}
	return h
}

func (h *ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(h.config.ToJSONString()))
}
