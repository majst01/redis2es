package main

import "strings"

// FilterStream passes the data between filters.
// future plugin api:
// https://medium.com/learning-the-go-programming-language/writing-modular-go-programs-with-plugins-ec46381ee1a9
type FilterStream struct {
	mapContent  map[string]interface{}
	jsonContent string
	indexName   string
}

type lowercase string

func (l lowercase) Name() string {
	return "lowercase keys filter"
}

func (l lowercase) Filter(stream *FilterStream) (*FilterStream, error) {
	for k, v := range stream.mapContent {
		lowerCaseKey := strings.ToLower(k)
		if lowerCaseKey == k {
			continue
		}
		stream.mapContent[lowerCaseKey] = v
		delete(stream.mapContent, k)
	}
	return stream, nil
}

// FilterPlugin exported symbol makes this plugin usable.
var FilterPlugin lowercase
