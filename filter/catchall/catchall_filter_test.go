package main

import (
	"testing"

	"github.com/majst01/redis2es/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	cf := catchallFilter{}

	stream := &filter.Stream{}
	err := cf.Filter(stream, "logstash")
	require.Nil(t, err, "error is expected to be nil")
	assert.Contains(t, stream.IndexName, "logstash-catchall-20", "catchall is expexted in the indexname")
}
