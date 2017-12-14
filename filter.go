package main

import (
	"path"
	"path/filepath"
	"plugin"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/json-iterator/go"
	"github.com/majst01/redis2es/filter"
)

// FilterPlugin modifies the input, which is a map representation of the json received
// to a output map or errors out.
type FilterPlugin interface {
	Name() string
	Filter(input *filter.Stream) error
}

const (
	filterDirectory = "lib"
	filterSuffix    = "_filter.so"
)

var (
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
	filters    []FilterPlugin
	filterGlob = path.Join(filterDirectory, "*"+filterSuffix)
)

func (r *redisClient) loadFilters() {
	log.WithFields(log.Fields{"init": "initialize filters"}).Info("loadfilters::")

	filters = []FilterPlugin{}

	// find all additional filter plugins and add them to a list
	files, err := filepath.Glob(filterGlob)
	if err != nil {
		log.WithFields(log.Fields{"cannot open filters": err}).Error("loadfilters:")
	}
	for _, file := range files {
		if !r.isFilterEnabled(file) {
			log.WithFields(log.Fields{"filter disabled": file}).Info("loadfilters:")
			continue
		}
		// load module
		// 1. open the so file to load the symbols
		plugin, err := plugin.Open(file)
		if err != nil {
			log.WithFields(log.Fields{"opening filter failed": file}).Error("loadfilters:")
		}
		// 2. look up a symbol (an exported function or variable)
		// in this case, variable FilterPlugin
		module, err := plugin.Lookup("FilterPlugin")
		if err != nil {
			log.WithFields(log.Fields{"FilterPlugin not detected": file}).Error("loadfilters:")
		}

		// 3. Assert that loaded symbol is of a desired type
		filter, ok := module.(FilterPlugin)
		if !ok {
			log.WithFields(log.Fields{"FilterPlugin interface not detected": file}).Error("loadfilters:")
		}
		log.WithFields(log.Fields{"filtername": filter.Name(), "filter shared lib": file}).Info("loadfilters:")
		filters = append(filters, filter)
	}
}

func (r *redisClient) isFilterEnabled(file string) bool {
	filtername := getFilterName(file)
	for _, filtered := range r.enabledFilters {
		if filtered == filtername {
			return true
		}
	}
	return false
}

func getFilterName(filename string) string {
	base := path.Base(filename)
	if !strings.HasSuffix(filename, filterSuffix) {
		return ""
	}
	filtername := strings.TrimSuffix(base, filterSuffix)
	return filtername
}

// getFilters is used to show which filters are available in total.
func getFilters() []string {
	var filters []string
	files, err := filepath.Glob(filterGlob)
	if err != nil {
		log.WithFields(log.Fields{"cannot open filters": err}).Error("filter:")
	}
	for _, file := range files {
		filtername := getFilterName(file)
		filters = append(filters, filtername)
	}
	return filters
}

func processFilter(input string) (*filter.Stream, error) {
	start := time.Now()
	stream := &filter.Stream{
		JSONContent: input,
	}
	if len(filters) == 0 {
		log.Debug("filter: no filters defined, returning original")
		return stream, nil
	}
	err := stream.Unmarshal()
	if err != nil {
		return stream, err
	}

	// check if contract in any case is present, lowercase then
	for _, f := range filters {
		s := time.Now()
		log.WithFields(log.Fields{"call filter:": f.Name()}).Debug("filter:")
		err = f.Filter(stream)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("filter:")
		}
		log.WithFields(log.Fields{"filtername": f.Name(), "duration": time.Now().Sub(s)}).Debug("filter:")
	}

	err = stream.Marshal()
	if err != nil {
		return stream, err
	}
	log.WithFields(log.Fields{"totalduration": time.Now().Sub(start)}).Debug("filter:")
	return stream, nil
}
