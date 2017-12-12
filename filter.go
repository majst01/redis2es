package main

import (
	"fmt"
	"path/filepath"
	"plugin"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/json-iterator/go"
)

// FilterStream passes the data between filters.
// future plugin api:
// https://medium.com/learning-the-go-programming-language/writing-modular-go-programs-with-plugins-ec46381ee1a9
type FilterStream struct {
	mapContent  map[string]interface{}
	jsonContent string
	indexName   string
}

// FilterPlugin modifies the input, which is a map representation of the json received
// to a output map or errors out.
type FilterPlugin interface {
	Name() string
	Filter(input *FilterStream) (*FilterStream, error)
}

var (
	json    = jsoniter.ConfigCompatibleWithStandardLibrary
	filters []FilterPlugin
)

func init() {
	log.WithFields(log.Fields{"init": "initialize filters"}).Info("filter:")

	filters = []FilterPlugin{}
	// add builtin filters
	filters = append(filters, ContractFilter{})

	// find all additional filter plugins and add them to a list
	files, err := filepath.Glob("./*_filter.so")
	if err != nil {
		log.WithFields(log.Fields{"cannot open filters": err}).Error("filter:")
	}
	for _, file := range files {
		// load module
		// 1. open the so file to load the symbols
		plugin, err := plugin.Open(file)
		if err != nil {
			log.WithFields(log.Fields{"opening filter failed": file}).Error("filter:")
		}
		// 2. look up a symbol (an exported function or variable)
		// in this case, variable FilterPlugin
		module, err := plugin.Lookup("FilterPlugin")
		if err != nil {
			log.WithFields(log.Fields{"FilterPlugin not detected": file}).Error("filter:")
		}

		// 3. Assert that loaded symbol is of a desired type
		filter, ok := module.(FilterPlugin)
		if !ok {
			log.WithFields(log.Fields{"FilterPlugin interface not detected": file}).Error("filter:")
		}
		log.WithFields(log.Fields{"filter shared lib": file, "filtername": filter.Name()}).Info("filter:")
		filters = append(filters, filter)
	}
}

func filter(input string) (*FilterStream, error) {
	start := time.Now()
	data := make(map[string]interface{})
	stream := &FilterStream{
		mapContent:  data,
		indexName:   fmt.Sprintf("logstash-catchall-%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day()),
		jsonContent: input,
	}
	if len(filters) == 0 {
		log.Debug("filter: no filters defined, returning original")
		return stream, nil
	}
	err := json.UnmarshalFromString(input, &data)
	if err != nil {
		return stream, fmt.Errorf("cannot decode data:%v", err)
	}

	// check if contract in any case is present, lowercase then
	for _, f := range filters {
		s := time.Now()
		log.WithFields(log.Fields{"call filter:": f.Name()}).Debug("filter:")
		stream, err = f.Filter(stream)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("filter:")
		}
		log.WithFields(log.Fields{"filtername": f.Name(), "duration": time.Now().Sub(s)}).Debug("filter:")
	}

	result, err := json.MarshalToString(&data)
	if err != nil {
		return stream, fmt.Errorf("cannot encode data:%v", err)
	}
	stream.jsonContent = result
	log.WithFields(log.Fields{"totalduration": time.Now().Sub(start)}).Debug("filter:")
	return stream, nil
}
