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

func getStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value != "" {
		v, err := strconv.Atoi(value)
		if err == nil {
			return v
		}
	}
	return defaultValue
}
func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value != "" {
		v, err := strconv.ParseBool(value)
		if err == nil {
			return v
		}
	}
	return defaultValue
}

func main() {
	redisKey := getStringEnv("REDIS_KEY", "logstash")
	redisHost := getStringEnv("REDIS_HOST", "127.0.0.1")
	redisPort := getIntEnv("REDIS_PORT", 6379)
	redisDB := getIntEnv("REDIS_DB", 0)
	redisPassword := getStringEnv("REDIS_PASSWORD", "")
	redisUseTLS := getBoolEnv("REDIS_USETLS", false)
	redisTLSSkipverify := getBoolEnv("REDIS_TLSSKIPVERIFY", false)

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
