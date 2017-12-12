package main

import (
	"fmt"
	"strings"
	"time"
)

type FilterStream struct {
	mapContent  map[string]interface{}
	jsonContent string
	indexName   string
}

// CustomerFilter check if customer in any case is present, lowercase it and and calculate indexname of it.
type CustomerFilter struct {
}

// Name required to be a FilterPlugin
func (cf CustomerFilter) Name() string {
	return "customerfilter"
}

// Filter required to be a FilterPlugin
func (cf CustomerFilter) Filter(stream *FilterStream) (*FilterStream, error) {
	for k, v := range stream.mapContent {
		if strings.ToLower(k) == "customer" {
			vString, ok := v.(string)
			if !ok {
				return stream, fmt.Errorf("customer is not a string")
			}
			oldValue := strings.ToLower(vString)
			delete(stream.mapContent, k)
			stream.mapContent["customer"] = oldValue
			stream.indexName = fmt.Sprintf("logstash-%s-%d.%d.%d", oldValue, time.Now().Year(), time.Now().Month(), time.Now().Day())
		}
	}

	return stream, nil
}

// FilterPlugin exported symbol makes this plugin usable.
var FilterPlugin CustomerFilter
