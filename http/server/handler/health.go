package handler

import (
	"github.com/heptiolabs/healthcheck"
	"github.com/rs/zerolog/log"
	"github.com/spals/starter-kit/http/server/config"
)

// NewHealthCheckHandler ...
// Creates live and readiness checks based off of HTTPServer configuration
//
// See https://github.com/heptiolabs/healthcheck/blob/master/README.md
func NewHealthCheckHandler(config *config.HTTPServerConfig) *healthcheck.Handler {
	healthCheckHandler := healthcheck.NewHandler()
	configureLivenessChecks(config, healthCheckHandler)
	configureReadinessChecks(config, healthCheckHandler)

	return &healthCheckHandler
}

// ========== Private Helpers ==========

func configureLivenessChecks(config *config.HTTPServerConfig, healthCheckHandler healthcheck.Handler) {
	log.Debug().Str("name", "goroutine-threshold").Int("check", config.LivenessConfig.MaxGoRoutines).Msg("Adding liveness check")
	healthCheckHandler.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(config.LivenessConfig.MaxGoRoutines))
}

func configureReadinessChecks(config *config.HTTPServerConfig, healthCheckHandler healthcheck.Handler) {

}
