package main

import (
	"strings"

	"github.com/majst01/redis2es/filter"
)

type lowercase string

func (l lowercase) Name() string {
	return "lowercase keys filter"
}

func (l lowercase) Filter(stream *filter.Stream) error {
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

// FilterPlugin exported symbol makes this plugin usable.
var FilterPlugin lowercase
