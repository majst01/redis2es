package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/olivere/elastic"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	redisPool := newPool("127.0.0.1", 6379, 0, "", false, false)
	client, err := elastic.NewSimpleClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		// Handle error
		panic(err)
	}
	defer client.Stop()

	rc := &redisClient{
		pool: redisPool,
		key:  "logstash",
		ec:   client,
	}

	err = rc.consume()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("main:")
	}

}
