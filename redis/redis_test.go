package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkConsume(b *testing.B) {
	rp := newPool("172.21.0.1", 6379, 0, "", false, false)

	client := rp.Get()
	for i := 0; i < 100000; i++ {
		_, err := client.Do("RPUSH", "logstash", "{\"UpperCaseKey\":\"value\",\"Contract\":\"FI-SP\"}")
		assert.Nil(b, err, "error must be nil")
	}
}
