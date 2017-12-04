package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/json-iterator/go"
)

// FilterStream passes the data between filters.
// future plugin api:
// https://medium.com/learning-the-go-programming-language/writing-modular-go-programs-with-plugins-ec46381ee1a9
type FilterStream struct {
	payload   map[string]interface{}
	json      string
	indexName string
}

// Filter modifies the input, which is a map representation of the json received
// to a output map or errors out.
type Filter interface {
	Name() string
	Filter(input *FilterStream) (*FilterStream, error)
}

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func filters() []Filter {
	f := []Filter{}
	f = append(f, ContractFilter{})
	return f
}

func filter(input string) (*FilterStream, error) {
	start := time.Now()
	data := make(map[string]interface{})
	err := json.UnmarshalFromString(input, &data)
	stream := &FilterStream{
		payload:   data,
		indexName: fmt.Sprintf("logstash-catchall-%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day()),
	}
	if err != nil {
		return stream, fmt.Errorf("cannot decode data:%v", err)
	}

	// check if contract in any case is present, lowercase then
	for _, f := range filters() {
		s := time.Now()
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
	stream.json = result
	log.WithFields(log.Fields{"totalduration": time.Now().Sub(start)}).Debug("filter:")
	return stream, nil
}
