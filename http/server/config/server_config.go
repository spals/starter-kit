package config

import (
	"context"
	"encoding/json"
	"fmt"
	nativelog "log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
)

// HTTPServerConfig ...
// Full configuration for HTTPServer
//
// See https://github.com/sethvargo/go-envconfig/blob/main/README.md
type HTTPServerConfig struct {
	Dev             bool          `env:"DEV,default=false"`
	Port            int           `env:"PORT,default=0"`
	LogLevel        string        `env:"LOG_LEVEL"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,default=1s"`

	LivenessConfig  *LivenessConfig  `env:",prefix=LIVENESS_"`
	ReadinessConfig *ReadinessConfig `env:"prefix=READINESS_"`
}

// NewHTTPServerConfig ...
func NewHTTPServerConfig(l envconfig.Lookuper) *HTTPServerConfig {
	ctx := context.Background()
	var config HTTPServerConfig

	if err := envconfig.ProcessWith(ctx, &config, l); err != nil {
		nativelog.Fatalf("HTTPServerConfig parse failure: %s", err)
		os.Exit(1)
	}

	log.Logger = makeLogger(&config)
	log.Debug().Interface("config", config).Msg("HTTPServerConfig parsed")
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

// ========== Private Helpers ==========

func makeLogger(config *HTTPServerConfig) zerolog.Logger {
	logLevel, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		nativelog.Fatalf("Error while parsing log level: %s. Available log levels are (trace|debug|info|warn|error|fatal|panic)", err)
	} else if logLevel == zerolog.NoLevel {
		nativelog.Fatalf("No log level configured. Please specify a log level (trace|debug|info|warn|error|fatal|panic)")
	}
	zerolog.SetGlobalLevel(logLevel)

	if config.Dev {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}

		return zerolog.New(output).With().Timestamp().Logger()
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		return zerolog.New(os.Stderr).With().Timestamp().Logger()
	}
}
