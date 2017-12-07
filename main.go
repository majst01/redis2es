package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/olivere/elastic"
)

// Specification of all configuration needed.
type Specification struct {
	Key           string `default:"logstash"`
	Host          string `default:"127.0.0.1"`
	Port          int    `default:"6379"`
	DB            int    `default:"0"`
	Password      string
	UseTLS        bool     `default:"false"`
	TLSSkipVerify bool     `default:"false"`
	ElasticURLs   []string `default:"http://127.0.0.1:9200"`
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	if os.Getenv("DEBUG") != "" {
		log.WithFields(log.Fields{"debug enabled": true}).Info("main:")
		log.SetLevel(log.DebugLevel)
	} else {
		log.WithFields(log.Fields{"debug disabled": true}).Info("main:")
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	var spec Specification
	err := envconfig.Process("redis2es", &spec)
	if err != nil {
		log.Fatal(err.Error())
	}

	redisPool := newPool(spec.Host, spec.Port, spec.DB, spec.Password, spec.UseTLS, spec.TLSSkipVerify)

	client, err := elastic.NewSimpleClient(elastic.SetURL(spec.ElasticURLs...))
	if err != nil {
		log.WithFields(log.Fields{"error connecting to elastic": err}).Error("main:")
	}
	defer client.Stop()

	rc := &redisClient{
		pool:    redisPool,
		key:     spec.Key,
		ec:      client,
		indexes: make(map[string]*elastic.BulkService),
	}

	err = rc.consume()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("main:")
	}

}
