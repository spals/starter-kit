package config

import (
	"context"
	"log"
	"time"

	"encoding/json"

	"github.com/sethvargo/go-envconfig"
)

// HTTPServerConfig ...
// See https://github.com/sethvargo/go-envconfig/blob/main/README.md
type HTTPServerConfig struct {
	Port            int16         `env:"PORT,default=8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,default=1s"`
}

// NewHTTPServerConfig ...
func NewHTTPServerConfig(l envconfig.Lookuper) *HTTPServerConfig {
	ctx := context.Background()
	var config HTTPServerConfig

	if err := envconfig.ProcessWith(ctx, &config, l); err != nil {
		log.Fatalf("HTTPServerConfig parse failure: %s", err)
	}

	return &config
}

// ToJSONString ...
func (c *HTTPServerConfig) ToJSONString() string {
	json, _ := json.Marshal(c)
	return string(json)
}
