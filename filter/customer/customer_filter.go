package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/majst01/redis2es/filter"
)

// CustomerFilter check if customer in any case is present, lowercase it and and calculate indexname of it.
type customerFilter struct {
}

const customer = "customer"

// Name required to be a filter.Plugin
func (cf customerFilter) Name() string {
	return "customerfilter"
}

// Filter required to be a filter.Plugin
func (cf customerFilter) Filter(stream *filter.Stream) error {
	for k, v := range stream.MapContent {
		if strings.ToLower(k) == customer {
			vString, ok := v.(string)
			if !ok {
				return fmt.Errorf("%s is not a string", customer)
			}
			oldValue := strings.ToLower(vString)
			delete(stream.MapContent, k)
			stream.MapContent[customer] = oldValue
			stream.IndexName = fmt.Sprintf("logstash-%s-%d.%d.%d", oldValue, time.Now().Year(), time.Now().Month(), time.Now().Day())
		}
	}

	return nil
}

// Plugin exported symbol makes this plugin usable.
var Plugin customerFilter
