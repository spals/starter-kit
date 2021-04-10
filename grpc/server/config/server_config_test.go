package config_test

import (
	"testing"

	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/grpc/server/config"
	"github.com/stretchr/testify/assert"
)

func TestBlankConfig(t *testing.T) {
	assert := assert.New(t)

	configMap := make(map[string]string)
	lookuper := envconfig.MapLookuper(configMap)

	config := config.NewGrpcServerConfig(lookuper)
	assert.Equal(int32(0), config.Port)
}

func TestPortLookup(t *testing.T) {
	assert := assert.New(t)

	configMap := make(map[string]string)
	configMap["PORT"] = "18080"

	lookuper := envconfig.MapLookuper(configMap)

	config := config.NewGrpcServerConfig(lookuper)
	assert.Equal(int32(18080), config.Port)
}
