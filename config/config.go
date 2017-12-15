package config

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Specification of all configuration needed.
type Specification struct {
	Key             string        `default:"logstash" desc:"the redis key where to listen to" required:"False"`
	Host            string        `default:"127.0.0.1" desc:"the redis server host/ip" required:"False"`
	Port            int           `default:"6379" desc:"the redis server port" required:"False"`
	DB              int           `default:"0" desc:"the redis database" required:"False"`
	Password        string        `desc:"the redis password" required:"False"`
	UseTLS          bool          `default:"false" desc:"connect to redis using tls" required:"False"`
	TLSSkipVerify   bool          `default:"false" desc:"if connection to redis via tls, skip tls certificate verification" required:"False"`
	ElasticURLs     []string      `default:"http://127.0.0.1:9200" desc:"the elasticsearch connection url, separated by comma for many es servers" required:"False"`
	ElasticUsername string        `desc:"username to connect to elasticsearch" required:"False"`
	ElasticPassword string        `desc:"password to connect to elasticsearch" required:"False"`
	BulkSize        int           `default:"1000" desc:"writes to elastic are done in bulks of bulkSize" required:"False"`
	BulkTicker      time.Duration `default:"2s" desc:"duration (go time.Duration format) between bulk writes to elastic" required:"False"`
	PoolSize        int           `default:"2" desc:"pool of workers to consume redis messages and write to elasticsearch" required:"False"`
	EnabledFilters  []string      `default:"catchall,customer" desc:"comma separated list of filters to be used, get a list of availablefilters with -l" required:"False"`
	StatsInterval   time.Duration `default:"60s" desc:"the interval on which bulkprocessor stats should be printed out" required:"False"`
	Debug           bool          `default:"false" desc:"turn on debug log" required:"False"`
}

// Log prints all config to log
func (s *Specification) Log() {
	log.WithFields(log.Fields{
		"key":             s.Key,
		"host":            s.Host,
		"port":            s.Port,
		"db":              s.DB,
		"password":        s.Password,
		"usetls":          s.UseTLS,
		"tlsskipverify":   s.TLSSkipVerify,
		"elasticurls":     s.ElasticURLs,
		"elasticusername": s.ElasticUsername,
		"elasticpassword": s.ElasticPassword,
		"bulksize":        s.BulkSize,
		"bulkticker":      s.BulkTicker,
		"poolsize":        s.PoolSize,
		"enabledfilters":  s.EnabledFilters,
		"statsinterval":   s.StatsInterval,
		"debug":           s.Debug,
	}).Info("config:")
}
