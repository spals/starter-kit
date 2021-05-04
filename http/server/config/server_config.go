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

	// Configured logger used for request logging
	ReqLogger zerolog.Logger
}

// NewHTTPServerConfig ...
func NewHTTPServerConfig(l envconfig.Lookuper) *HTTPServerConfig {
	ctx := context.Background()
	var config HTTPServerConfig

	if err := envconfig.ProcessWith(ctx, &config, l); err != nil {
		nativelog.Fatalf("HTTPServerConfig parse failure: %s", err)
		os.Exit(1)
	}

	// Configure logging as early as possible (i.e. as soon as we have a parsed configuration)
	config.configureLogging()
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

func (c *HTTPServerConfig) configureLogging() {
	// Set the default logger as the application logger
	log.Logger = c.newLogger().With().Str("system", "starter-kit-http").Logger()
	c.ReqLogger = c.newLogger().With().Str("system", "http-request").Logger()
}

func (c *HTTPServerConfig) newLogger() zerolog.Logger {
	logLevel, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		nativelog.Fatalf("Error while parsing log level: %s. Available log levels are (trace|debug|info|warn|error|fatal|panic)", err)
	} else if logLevel == zerolog.NoLevel {
		nativelog.Fatalf("No log level configured. Please specify a log level (trace|debug|info|warn|error|fatal|panic)")
	}
	zerolog.SetGlobalLevel(logLevel)

	if c.Dev {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("[%s]:", i)
		}

		return zerolog.New(output).With().Timestamp().Caller().Logger()
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		return zerolog.New(os.Stderr).With().Timestamp().Logger()
	}
}
