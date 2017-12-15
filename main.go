package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

var (
	// EnabledFilters specifies which filters to use.
	EnabledFilters []string
)

type document struct {
	indexName string
	body      string
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	var spec Specification
	envconfig.MustProcess("redis2es", &spec)
	spec.log()
	if len(os.Args) > 1 {
		if os.Args[1] == "-l" {
			log.WithFields(log.Fields{"filters": getFilters()}).Info("main:")
			os.Exit(0)
		}
		envconfig.Usage("redis2es", &spec)
		os.Exit(1)
	}

	if spec.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	redisPool := newPool(spec.Host, spec.Port, spec.DB, spec.Password, spec.UseTLS, spec.TLSSkipVerify)

	ec := NewElasticClient(spec)
	defer ec.close()

	rc := &redisClient{
		pool:           redisPool,
		key:            spec.Key,
		enabledFilters: spec.EnabledFilters,
	}

	rc.loadFilters()

	for i := 0; i < spec.PoolSize; i++ {
		documents := make(chan document, 10)
		go ec.index(documents)
		go rc.consume(documents)
	}

	// Stay in forground
	ec.stats()
}
