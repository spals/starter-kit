package config_test

import (
	"fmt"
	"testing"

	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/grpc/server/config"
	"github.com/stretchr/testify/assert"
)

func TestNewGrpcServerConfig(t *testing.T) {
	assert := assert.New(t)

	configMap := make(map[string]string)
	configMap["PORT"] = fmt.Sprintf("%d", 18080)

	lookuper := envconfig.MapLookuper(configMap)

	config := config.NewGrpcServerConfig(lookuper)
	assert.Equal(18080, config.Port)
}
