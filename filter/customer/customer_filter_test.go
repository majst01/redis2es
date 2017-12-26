package main

import (
	"testing"

	"github.com/majst01/redis2es/filter"
	"github.com/stretchr/testify/assert"
)

func TestCustomerFilter(t *testing.T) {
	cf := customerFilter{}

	stream := &filter.Stream{
		JSONContent: "{\"key\":\"value\"}",
	}
	err := stream.Unmarshal()

	err = cf.Filter(stream, "logstash")
	assert.Nil(t, err, "no error is expected")
	assert.Equal(t, "value", stream.MapContent["key"], "output may never be empty")

	stream.JSONContent = "{\"key\":\"value\", \"Customer\":\"TestCustomer\"}"
	err = stream.Unmarshal()
	assert.Nil(t, err, "no error is expected")
	err = cf.Filter(stream, "logstash")
	assert.Nil(t, err, "no error is expected")
	assert.Contains(t, stream.IndexName, "testcustomer", "index must contain testcustomer")
	assert.Equal(t, "testcustomer", stream.MapContent["customer"], "customer expected")
}
