package main

import (
	"context"
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/olivere/elastic"
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

	client, err := elastic.NewSimpleClient(elastic.SetURL(spec.ElasticURLs...))
	if err != nil {
		log.WithFields(log.Fields{"error connecting to elastic": err}).Error("main:")
	}
	bulk, err := client.BulkProcessor().
		Name("BackgroundWorker-1").
		Workers(spec.PoolSize).         // number of workers
		BulkActions(spec.BulkSize).     // commit if # requests >= BulkSize
		BulkSize(2 << 20).              // commit if size of requests >= 2 MB
		FlushInterval(spec.BulkTicker). // commit every given interval
		Stats(true).                    // collect stats
		Do(context.Background())
	if err != nil {
		log.WithFields(log.Fields{"error creating bulkprocessor": err}).Fatal("main:")
	}

	defer bulk.Close()
	defer client.Stop()

	rc := &redisClient{
		pool:           redisPool,
		key:            spec.Key,
		ec:             client,
		bulkProcessor:  bulk,
		enabledFilters: spec.EnabledFilters,
	}

	rc.loadFilters()

	for i := 0; i < spec.PoolSize; i++ {
		documents := make(chan document, 10)
		go rc.index(documents)
		go rc.consume(documents)
	}

	// Stay in forground
	rc.stats()
}
