package config

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Redis specific configuration
type Redis struct {
	Key           string `default:"logstash" desc:"the redis key where to listen to" required:"False"`
	Host          string `default:"127.0.0.1" desc:"the redis server host/ip" required:"False"`
	Port          int    `default:"6379" desc:"the redis server port" required:"False"`
	DB            int    `default:"0" desc:"the redis database" required:"False"`
	Password      string `desc:"the redis password" required:"False"`
	UseTLS        bool   `default:"false" desc:"connect to redis using tls" required:"False"`
	TLSSkipVerify bool   `default:"false" desc:"if connection to redis via tls, skip tls certificate verification" required:"False"`
	PoolSize      int    `default:"2" desc:"pool of workers to consume redis messages" required:"False"`
}

// Elastic specific configuration
type Elastic struct {
	URLs           []string      `default:"http://127.0.0.1:9200" desc:"the elasticsearch connection url, separated by comma for many es servers" required:"False"`
	Username       string        `desc:"username to connect to elasticsearch" required:"False"`
	Password       string        `desc:"password to connect to elasticsearch" required:"False"`
	BulkSize       int           `default:"1000" desc:"writes to elastic are done in bulks of bulkSize" required:"False"`
	BulkTicker     time.Duration `default:"2s" desc:"duration (go time.Duration format) between bulk writes to elastic" required:"False"`
	PoolSize       int           `default:"2" desc:"pool of workers to write to elasticsearch" required:"False"`
	EnabledFilters []string      `default:"catchall,customer" desc:"comma separated list of filters to be used, get a list of availablefilters with -l" required:"False"`
	StatsInterval  time.Duration `default:"60s" desc:"the interval on which bulkprocessor stats should be printed out" required:"False"`
	IndexPrefix    string        `default:"logstash" desc:"the prefix to be used for indexes" required:"False"`
}

// Specification of all configuration needed.
type Specification struct {
	Redis   Redis
	Elastic Elastic
	Debug   bool `default:"false" desc:"turn on debug log" required:"False"`
}

// Log prints all config to log
func (s *Specification) Log() {
	log.WithFields(log.Fields{
		"redis.key":              s.Redis.Key,
		"redis.host":             s.Redis.Host,
		"redis.port":             s.Redis.Port,
		"redis.db":               s.Redis.DB,
		"redis.password":         s.Redis.Password,
		"redis.usetls":           s.Redis.UseTLS,
		"redis.tlsskipverify":    s.Redis.TLSSkipVerify,
		"redis.poolsize":         s.Redis.PoolSize,
		"elastic.urls":           s.Elastic.URLs,
		"elastic.username":       s.Elastic.Username,
		"elastic.password":       s.Elastic.Password,
		"elastic.bulksize":       s.Elastic.BulkSize,
		"elastic.bulkticker":     s.Elastic.BulkTicker,
		"elastic.poolsize":       s.Elastic.PoolSize,
		"elastic.enabledfilters": s.Elastic.EnabledFilters,
		"elastic.statsinterval":  s.Elastic.StatsInterval,
		"elastic.indexprefix":    s.Elastic.IndexPrefix,
		"debug":                  s.Debug,
	}).Info("config:")
}
