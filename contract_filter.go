package main

import (
	"fmt"
	"strings"
	"time"
)

type ContractFilter struct {
}

func (cf ContractFilter) Name() string {
	return "contractfilter"
}

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
