package config

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

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

	LivenessConfig  *LivenessConfig  `env:",prefix=LIVENESS_"`
	ReadinessConfig *ReadinessConfig `env:"prefix=READINESS_"`
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

	log.Printf("HTTPServerConfig parsed as \n%s", config.ToJSONString(true /*prettyPrint*/))
	return &config
}

// ToJSONString ...
func (c *HTTPServerConfig) ToJSONString(prettyPrint bool) string {
	if prettyPrint {
		json, _ := json.MarshalIndent(c, "", "  ")
		return string(json)
	}

	json, _ := json.Marshal(c)
	return string(json)
}
