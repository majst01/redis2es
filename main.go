package main

import (
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/olivere/elastic"
)

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
	redisKey := "logstash"
	if os.Getenv("REDIS_KEY") != "" {
		redisKey = os.Getenv("REDIS_KEY")
	}
	redisHost := "127.0.0.1"
	if os.Getenv("REDIS_HOST") != "" {
		redisHost = os.Getenv("REDIS_HOST")
	}
	redisPort := 6379
	if os.Getenv("REDIS_Port") != "" {
		port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
		if err != nil {
			log.WithFields(log.Fields{"REDIS_PORT is not numeric": err}).Error("main:")
		}
		redisPort = port
	}
	redisDB := 0
	if os.Getenv("REDIS_DB") != "" {
		db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			log.WithFields(log.Fields{"REDIS_DB is not numeric": err}).Error("main:")
		}
		redisDB = db
	}
	redisPassword := ""
	if os.Getenv("REDIS_PASSWORD") != "" {
		redisPassword = os.Getenv("REDIS_PASSWORD")
	}
	redisUseTLS := false
	if os.Getenv("REDIS_USETLS") != "" {
		useTLS, err := strconv.ParseBool(os.Getenv("REDIS_USETLS"))
		if err != nil {
			log.WithFields(log.Fields{"REDIS_USETLS is not a bool": err}).Error("main:")
		}
		redisUseTLS = useTLS
	}
	redisTLSSkipverify := false
	if os.Getenv("REDIS_TLSSKIPVERIFY") != "" {
		skiverify, err := strconv.ParseBool(os.Getenv("REDIS_TLSSKIPVERIFY"))
		if err != nil {
			log.WithFields(log.Fields{"REDIS_TLSSKIPVERIFY is not a bool": err}).Error("main:")
		}
		redisTLSSkipverify = skiverify
	}

	redisPool := newPool(redisHost, redisPort, redisDB, redisPassword, redisUseTLS, redisTLSSkipverify)

	elasticURLs := []string{"http://127.0.0.1:9200"}
	if os.Getenv("ELASTIC_URLS") != "" {
		elasticURLs = strings.Split(os.Getenv("ELASTIC_URLS"), ",")
	}
	client, err := elastic.NewSimpleClient(elastic.SetURL(elasticURLs...))
	if err != nil {
		log.WithFields(log.Fields{"error connecting to elastic": err}).Error("main:")
	}
	defer client.Stop()

	rc := &redisClient{
		pool: redisPool,
		key:  redisKey,
		ec:   client,
	}

	err = rc.consume()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("main:")
	}

}
