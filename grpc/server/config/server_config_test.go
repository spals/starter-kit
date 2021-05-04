package config_test

import (
	"testing"

	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/grpc/server/config"
	"github.com/stretchr/testify/assert"
)

func TestDevLookup(t *testing.T) {
	assert := assert.New(t)

	configMap := make(map[string]string)
	configMap["LOG_LEVEL"] = "trace"
	configMap["DEV"] = "true"

	lookuper := envconfig.MapLookuper(configMap)

	config := config.NewGrpcServerConfig(lookuper)
	assert.Equal(true, config.Dev)
}

func TestLogLevelLookup(t *testing.T) {
	assert := assert.New(t)

	configMap := make(map[string]string)
	configMap["LOG_LEVEL"] = "info"

	lookuper := envconfig.MapLookuper(configMap)

	config := config.NewGrpcServerConfig(lookuper)
	assert.Equal("info", config.LogLevel)
	assert.Equal(int32(0), config.Port)
}

func TestPortLookup(t *testing.T) {
	assert := assert.New(t)

	configMap := make(map[string]string)
	configMap["LOG_LEVEL"] = "trace"
	configMap["PORT"] = "18080"

	lookuper := envconfig.MapLookuper(configMap)

	config := config.NewGrpcServerConfig(lookuper)
	assert.Equal(int32(18080), config.Port)
}
