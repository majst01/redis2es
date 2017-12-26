package main

import (
	"testing"

	"github.com/majst01/redis2es/filter"
	"github.com/stretchr/testify/assert"
)

func TestLowerCaseFilter(t *testing.T) {
	lc := &lowercaseFilter{}

	stream := &filter.Stream{
		JSONContent: "{\"Key\":\"value\"}",
	}
	err := stream.Unmarshal()

	err = lc.Filter(stream, "logstash")

	assert.Nil(t, err, "no error is expected")
	assert.Equal(t, "value", stream.MapContent["key"], "output may never be empty")

}
