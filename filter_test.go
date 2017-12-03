package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	input := "{\"key\":\"value\"}"
	output, contract, err := filter(input)
	assert.Nil(t, err, "no error is expected")
	assert.Equal(t, "catchall", contract, "contract may never be empty")
	assert.NotEqual(t, "\"key\":\"value\",\"contract\":\"catchall\"", output, "output may never be empty")

	input = "{\"key\":\"value\", \"Contract\":\"TestContract\"}"
	output, contract, err = filter(input)
	assert.Nil(t, err, "no error is expected")
	assert.Equal(t, "testcontract", contract)
	assert.NotEqual(t, "\"key\":\"value\",\"contract\":\"testcontract\"", output, "output may never be empty")

}

func BenchmarkFilter(b *testing.B) {
	input := "{\"key\":\"value\", \"Contract\":\"TestContract\"}"
	for i := 0; i < b.N; i++ {
		_, _, err := filter(input)
		if err != nil {
			assert.Fail(b, "%v", err)
		}
	}
}
