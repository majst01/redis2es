package main

import (
	"fmt"
	"strings"
	"time"
)

// ContractFilter check if contract in any case is present, lowercase it and and calculate indexname of it.
type ContractFilter struct {
}

// Name required to be a FilterPlugin
func (cf ContractFilter) Name() string {
	return "contractfilter"
}

// Filter required to be a FilterPlugin
func (cf ContractFilter) Filter(stream *FilterStream) (*FilterStream, error) {
	for k, v := range stream.payload {
		if strings.ToLower(k) == "contract" {
			vString, ok := v.(string)
			if !ok {
				return stream, fmt.Errorf("contract is not a string")
			}
			oldValue := strings.ToLower(vString)
			delete(stream.payload, k)
			stream.payload["contract"] = oldValue
			stream.indexName = fmt.Sprintf("logstash-%s-%d-%d-%d", oldValue, time.Now().Year(), time.Now().Month(), time.Now().Day())
		}
	}

	return stream, nil
}
