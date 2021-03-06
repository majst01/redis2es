package elastic

import (
	"fmt"
	"path"
	"path/filepath"
	"plugin"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/json-iterator/go"
	"github.com/majst01/redis2es/filter"
)

const (
	filterDirectory = "lib"
	filterSuffix    = "_filter.so"
)

var (
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
	filterGlob = path.Join(filterDirectory, "*"+filterSuffix)
)

func (e *Client) loadFilters() {
	log.WithFields(log.Fields{"init": "initialize filters"}).Info("loadfilters:")

	filters := []filter.Plugin{}

	// find all additional filter plugins and add them to a list
	files, err := filepath.Glob(filterGlob)
	if err != nil {
		log.WithFields(log.Fields{"cannot open filters": err}).Error("loadfilters:")
	}
	for _, file := range files {
		if !e.isFilterEnabled(file) {
			log.WithFields(log.Fields{"filter disabled": file}).Info("loadfilters:")
			continue
		}
		filter, err := loadFilter(file)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Info("loadfilters:")
			continue
		}
		filters = append(filters, filter)
	}
	e.filters = filters
}

func loadFilter(file string) (filter.Plugin, error) {
	// load module
	// 1. open the so file to load the symbols
	plugin, err := plugin.Open(file)
	if err != nil {
		return nil, fmt.Errorf("opening filter file %s failed with: %v", file, err)
	}
	// 2. look up a symbol (an exported function or variable)
	// in this case, variable Plugin
	module, err := plugin.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("Plugin not detected in %s with: %v", file, err)
	}

	// 3. Assert that loaded symbol is of a desired type
	filter, ok := module.(filter.Plugin)
	if !ok {
		return nil, fmt.Errorf("Plugin interface not detected in %s", file)
	}
	log.WithFields(log.Fields{"filtername": filter.Name(), "filterfile": file}).Info("loadfilters:")

	return filter, nil
}

// isFilterEnabled returns true if this filter is enabled by config.
func (e *Client) isFilterEnabled(file string) bool {
	filtername := getFilterName(file)
	for _, filtered := range e.enabledFilters {
		if filtered == filtername {
			return true
		}
	}
	return false
}

// getFilterName extracts a short filtername from filter file path.
func getFilterName(filename string) string {
	base := path.Base(filename)
	if !strings.HasSuffix(filename, filterSuffix) {
		return ""
	}
	filtername := strings.TrimSuffix(base, filterSuffix)
	return filtername
}

// GetFilters is used to show which filters are available in total.
func GetFilters() []string {
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

// ProcessFilter apply all filters enabled to the stream
func (e *Client) ProcessFilter(input string) (*filter.Stream, error) {
	start := time.Now()
	stream := &filter.Stream{
		JSONContent: input,
	}
	if len(e.filters) == 0 {
		log.Debug("filter: no filters defined, returning original")
		return stream, nil
	}
	err := stream.Unmarshal()
	if err != nil {
		return stream, err
	}

	for _, f := range e.filters {
		s := time.Now()
		log.WithFields(log.Fields{"call filter:": f.Name()}).Debug("filter:")
		err = f.Filter(stream, e.indexPrefix)
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
