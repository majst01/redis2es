package main

import (
	"os"
	"testing"

	"github.com/majst01/redis2es/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Equal(t, output.IndexName, "", "index must contain catchall")
	assert.Equal(t, "value", output.MapContent["key"], "expected to have a map representation of json input")

}

func TestGetFilterName(t *testing.T) {
	name := getFilterName("lib/test_filter.so")
	assert.Equal(t, name, "test", "filtername expected was not met.")

	name = getFilterName("lib/test_filter.txt")
	assert.Equal(t, name, "", "filtername expected was not met.")
}

func TestIsFilterEnabled(t *testing.T) {
	r := &redisClient{
		enabledFilters: []string{"noop", "dump"},
	}

	enabled := r.isFilterEnabled("lib/test_filter.so")
	assert.False(t, enabled, "filter is expected to disabled")

	enabled = r.isFilterEnabled("lib/noop_filter.so")
	assert.True(t, enabled, "filter is expected to disabled")

}

func TestGetFilters(t *testing.T) {
	defer os.RemoveAll("lib")
	err := os.RemoveAll("lib")
	err = os.Mkdir("lib", 0755)
	require.Nil(t, err)
	os.OpenFile("lib/test_filter.so", os.O_RDWR|os.O_CREATE, 0755)

	filters := getFilters()
	assert.True(t, len(filters) == 1, "one filter is expected")
	assert.Equal(t, "test", filters[0], "testfilter must be present")

	os.OpenFile("lib/noop_filter.so", os.O_RDWR|os.O_CREATE, 0755)
	filters = getFilters()
	assert.True(t, len(filters) == 2, "two filter is expected")
	assert.Equal(t, "noop", filters[0], "noopfilter must be present")
	assert.Equal(t, "test", filters[1], "testfilter must be present")
}

func TestLoadFilters(t *testing.T) {
	defer os.RemoveAll("lib")
	err := os.RemoveAll("lib")
	err = os.Mkdir("lib", 0755)
	require.Nil(t, err)
	os.OpenFile("lib/test_filter.so", os.O_RDWR|os.O_CREATE, 0755)

	// FIXME implement
	//r := &redisClient{
	//	enabledFilters: []string{"test"},
	//}

	// r.loadFilters()
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
