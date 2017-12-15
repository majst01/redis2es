package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
)

// RedisClient is used to read from redis
type RedisClient struct {
	key            string
	pool           *redis.Pool
	enabledFilters []string
	filters        []FilterPlugin
}

// NewRedisClient create a new instance of a redisClient
func NewRedisClient(spec Specification) *RedisClient {
	redisPool := newPool(spec.Host, spec.Port, spec.DB, spec.Password, spec.UseTLS, spec.TLSSkipVerify)

	rc := &RedisClient{
		pool:           redisPool,
		key:            spec.Key,
		enabledFilters: spec.EnabledFilters,
	}

	rc.loadFilters()
	return rc
}

func newPool(host string, port int, db int, password string, usetls, tlsskipverify bool) *redis.Pool {
	var c redis.Conn
	var err error
	server := fmt.Sprintf("%s:%d", host, port)
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			if usetls {
				c, err = redis.Dial("tcp", server, redis.DialDatabase(db),
					redis.DialUseTLS(usetls),
					redis.DialTLSSkipVerify(tlsskipverify),
				)
			} else {
				c, err = redis.Dial("tcp", server, redis.DialDatabase(db))
			}
			if err != nil {
				return nil, err
			}
			// In case redis needs authentication
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (r *RedisClient) readBlocking() (string, error) {
	c := r.pool.Get()
	defer c.Close()

	result, err := redis.Strings(c.Do("BLPOP", r.key, 0))
	if err != nil {
		time.Sleep(time.Duration(time.Millisecond * 200))
		return "", fmt.Errorf("error getting value to key: %s with: %v", r.key, err)
	}
	if len(result) != 2 {
		return "", fmt.Errorf("error getting value to key: %s, expected 2 entries, got %d", r.key, len(result))
	}
	return result[1], nil
}

func (r *RedisClient) pending() int {
	c := r.pool.Get()
	defer c.Close()
	v, err := redis.Int(c.Do("LLEN", r.key))
	if err != nil {
		return 0
	}
	return v
}

func (r *RedisClient) consume(documents chan document) {
	c := r.pool.Get()
	defer c.Close()

	for {
		result, err := r.readBlocking()
		if err != nil {
			log.WithFields(log.Fields{"error from BLPOP": err}).Error("consume:")
			continue
		}
		filtered, err := r.processFilter(result)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("consume:")
			continue
		}
		doc := document{
			indexName: filtered.IndexName,
			body:      filtered.JSONContent,
		}
		documents <- doc
	}
}
