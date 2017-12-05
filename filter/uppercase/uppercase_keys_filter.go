package main

import "strings"

// FilterStream passes the data between filters.
// future plugin api:
// https://medium.com/learning-the-go-programming-language/writing-modular-go-programs-with-plugins-ec46381ee1a9
type FilterStream struct {
	payload   map[string]interface{}
	json      string
	indexName string
}

type uppercase string

func (u uppercase) Name() string {
	return "uppercase keys filter"
}

func (u uppercase) Filter(stream *FilterStream) (*FilterStream, error) {
	for k, v := range stream.payload {
		lowerCaseKey := strings.ToLower(k)
		if lowerCaseKey == k {
			continue
		}
		stream.payload[lowerCaseKey] = v
		delete(stream.payload, k)
	}
	return stream, nil
}

var FilterPlugin uppercase
