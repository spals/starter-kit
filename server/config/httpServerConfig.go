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
// Full configuration for HTTPServer
//
// See https://github.com/sethvargo/go-envconfig/blob/main/README.md
type HTTPServerConfig struct {
	AssignRandomPort bool          `env:"ASSIGN_RANDOM_PORT,default=false"`
	Port             int           `env:"PORT,default=8080"`
	ShutdownTimeout  time.Duration `env:"SHUTDOWN_TIMEOUT,default=1s"`

	LivenessConfig *LivenessConfig `env:",prefix=LIVENESS_"`
}

// LivenessConfig ...
// Configuration used to check liveness for HTTPServer
type LivenessConfig struct {
	MaxGoRoutines int `env:"MAX_GO_ROUTINES,default=100"`
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
