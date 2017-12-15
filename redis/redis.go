package redis

import (
	"fmt"
	"time"

	"github.com/majst01/redis2es/config"
	"github.com/majst01/redis2es/elastic"
	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
)

// Client is used to read from redis
type Client struct {
	key  string
	pool *redis.Pool
	ec   *elastic.Client
}

// New create a new instance of a redis Client
func New(spec config.Redis, ec *elastic.Client) *Client {
	redisPool := newPool(spec.Host, spec.Port, spec.DB, spec.Password, spec.UseTLS, spec.TLSSkipVerify)

	rc := &Client{
		pool: redisPool,
		key:  spec.Key,
		ec:   ec,
	}

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

func (r *Client) readBlocking() (string, error) {
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

func (r *Client) pending() int {
	c := r.pool.Get()
	defer c.Close()
	v, err := redis.Int(c.Do("LLEN", r.key))
	if err != nil {
		return 0
	}
	return v
}

// Consume logs from redis
func (r *Client) Consume(documents chan elastic.Document) {
	c := r.pool.Get()
	defer c.Close()

	for {
		result, err := r.readBlocking()
		if err != nil {
			log.WithFields(log.Fields{"error from BLPOP": err}).Error("consume:")
			continue
		}
		filtered, err := r.ec.ProcessFilter(result)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("consume:")
			continue
		}
		doc := elastic.Document{
			IndexName: filtered.IndexName,
			Body:      filtered.JSONContent,
		}
		documents <- doc
	}
}
