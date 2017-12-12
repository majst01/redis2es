package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Specification of all configuration needed.
type Specification struct {
	Key            string        `default:"logstash" desc:"the redis key where to listen to" required:"False"`
	Host           string        `default:"127.0.0.1" desc:"the redis server host/ip" required:"False"`
	Port           int           `default:"6379" desc:"the redis server port" required:"False"`
	DB             int           `default:"0" desc:"the redis database" required:"False"`
	Password       string        `desc:"the redis password" required:"False"`
	UseTLS         bool          `default:"false" desc:"connect to redis using tls" required:"False"`
	TLSSkipVerify  bool          `default:"false" desc:"if connection to redis via tls, skip tls certificate verification" required:"False"`
	ElasticURLs    []string      `default:"http://127.0.0.1:9200" desc:"the elasticsearch connection url, seperated by comma for many es servers" required:"False"`
	BulkSize       int           `default:"1000" desc:"writes to elastic are done in bulks of bulkSize" required:"False"`
	BulkTicker     time.Duration `default:"2s" desc:"duration (go time.Duration format) between bulk writes to elastic" required:"False"`
	PoolSize       int           `default:"2" desc:"pool of workers to consume redis messages and write to elasticsearch" required:"False"`
	EnabledFilters []string      `default:"customer" desc:"comma seperated list of filters to be used, get a list of available filters with -l" required:"False"`
	Debug          bool          `default:"false" desc:"turn on debug log" required:"False"`
}

func (s *Specification) log() {
	log.WithFields(log.Fields{
		"key":            s.Key,
		"host":           s.Host,
		"port":           s.Port,
		"db":             s.DB,
		"password":       s.Password,
		"usetls":         s.UseTLS,
		"tlsskipverify":  s.TLSSkipVerify,
		"elasticurls":    s.ElasticURLs,
		"bulksize":       s.BulkSize,
		"bulkticker":     s.BulkTicker,
		"poolsize":       s.PoolSize,
		"enabledfilters": s.EnabledFilters,
		"debug":          s.Debug,
	}).Info("config:")
}
