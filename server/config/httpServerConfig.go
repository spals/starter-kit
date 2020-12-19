package config

import (
	"context"
	"log"
	"os"
	"time"

	"encoding/json"

	"github.com/sethvargo/go-envconfig"
)

// HTTPServerConfig ...
// See https://github.com/sethvargo/go-envconfig/blob/main/README.md
type HTTPServerConfig struct {
	Port            int           `env:"PORT,default=8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,default=1s"`
	UseRandomPort   bool          `env:"USE_RANDOM_PORT,default=false"`
}

// NewHTTPServerConfig ...
func NewHTTPServerConfig(l envconfig.Lookuper) *HTTPServerConfig {
	log.Print("Parsing HTTPServerConfig")
	ctx := context.Background()
	var config HTTPServerConfig

	if err := envconfig.ProcessWith(ctx, &config, l); err != nil {
		log.Fatalf("HTTPServerConfig parse failure: %s", err)
		os.Exit(1)
	}

	log.Printf("HTTPServerConfig parsed as %s", config.ToJSONString())
	return &config
}

// ToJSONString ...
func (c *HTTPServerConfig) ToJSONString() string {
	json, _ := json.Marshal(c)
	return string(json)
}
