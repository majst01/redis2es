package main

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/olivere/elastic"
)

type document struct {
	indexName string
	body      string
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

func main() {
	var spec Specification
	envconfig.MustProcess("redis2es", &spec)
	spec.log()
	if len(os.Args) > 1 {
		envconfig.Usage("redis2es", &spec)
		os.Exit(1)
	}
	if spec.Debug {
		log.WithFields(log.Fields{"debug enabled": true}).Info("main:")
		log.SetLevel(log.DebugLevel)
	} else {
		log.WithFields(log.Fields{"debug disabled": true}).Info("main:")
		log.SetLevel(log.InfoLevel)
	}

	redisPool := newPool(spec.Host, spec.Port, spec.DB, spec.Password, spec.UseTLS, spec.TLSSkipVerify)

	client, err := elastic.NewSimpleClient(elastic.SetURL(spec.ElasticURLs...))
	if err != nil {
		log.WithFields(log.Fields{"error connecting to elastic": err}).Error("main:")
	}
	defer client.Stop()

	rc := &redisClient{
		pool:     redisPool,
		key:      spec.Key,
		ec:       client,
		indexes:  make(map[string]*elastic.BulkService),
		bulkSize: spec.BulkSize,
	}

	// FIXME concurency configurable
	for i := 0; i < spec.PoolSize; i++ {
		documents := make(chan document, 10)
		go rc.index(documents)
		go rc.consume(documents)
	}
	// Stay in foreground
	fmt.Scanln()
}
