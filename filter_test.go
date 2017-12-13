package main

import (
	"testing"

	"github.com/majst01/redis2es/filter"
	"github.com/stretchr/testify/assert"
)

type noopfilter struct {
}

func (n noopfilter) Name() string {
	return "NOOP"
}
func (n noopfilter) Filter(stream *filter.Stream) error {
	return nil
}

func TestFilter(t *testing.T) {
	input := "{\"key\":\"value\"}"
	filters = []FilterPlugin{}
	filters = append(filters, noopfilter{})
	output, err := processFilter(input)
	assert.Nil(t, err, "no error is expected")
	assert.Contains(t, output.IndexName, "catchall", "index must contain catchall")
	assert.Equal(t, "value", output.MapContent["key"], "expected to have a map representation of json input")

}

func BenchmarkFilter(b *testing.B) {
	input := "{\"key\":\"value\", \"Contract\":\"TestContract\"}"
	for i := 0; i < b.N; i++ {
		_, err := processFilter(input)
		if err != nil {
			assert.Fail(b, "%v", err)
		}
	}
}
