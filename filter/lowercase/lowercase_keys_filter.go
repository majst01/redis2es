package main

import (
	"strings"

	"github.com/majst01/redis2es/filter"
)

type lowercaseFilter struct {
}

func (l lowercaseFilter) Name() string {
	return "lowercase keys filter"
}

func (l lowercaseFilter) Filter(stream *filter.Stream) error {
	for k, v := range stream.MapContent {
		lowerCaseKey := strings.ToLower(k)
		if lowerCaseKey == k {
			continue
		}
		stream.MapContent[lowerCaseKey] = v
		delete(stream.MapContent, k)
	}
	return nil
}

// Plugin exported symbol makes this plugin usable.
var Plugin lowercaseFilter
