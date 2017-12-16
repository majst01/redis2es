package main

import (
	"os"

	"github.com/majst01/redis2es/config"
	"github.com/majst01/redis2es/elastic"
	"github.com/majst01/redis2es/filter"
	"github.com/majst01/redis2es/redis"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	var spec config.Specification
	envconfig.MustProcess("redis2es", &spec)
	if len(os.Args) > 1 {
		if os.Args[1] == "-l" {
			log.WithFields(log.Fields{"filters": elastic.GetFilters()}).Info("main:")
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
	spec.Log()

	ec := elastic.New(spec.Elastic)
	defer ec.Close()

	rc := redis.New(spec.Redis, ec)

	stream := make(chan *filter.Stream, 10)
	go ec.Index(stream)
	for i := 0; i < spec.Redis.PoolSize; i++ {
		go rc.Consume(stream)
	}

	// Stay in forground
	ec.Stats(spec.Elastic.StatsInterval)
}
