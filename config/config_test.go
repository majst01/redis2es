package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kelseyhightower/envconfig"
)

func TestSpecification(t *testing.T) {
	var spec Specification

	envconfig.MustProcess("redis2es", &spec)

	assert.Equal(t, 2, spec.Redis.PoolSize, "redis.poolsize is expected to be 2")
}
