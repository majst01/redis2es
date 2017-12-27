package main

import (
	"fmt"
	"time"

	"github.com/majst01/redis2es/filter"
)

type catchallFilter struct{}

func (c catchallFilter) Name() string {
	return "catchall index filter"
}

func (c catchallFilter) Filter(stream *filter.Stream, indexPrefix string) error {
	stream.IndexName = fmt.Sprintf("%s-catchall-%d-%d-%d", indexPrefix, time.Now().Year(), time.Now().Month(), time.Now().Day())
	return nil
}

// Plugin exported symbol to make it a usable plugin
var Plugin catchallFilter
