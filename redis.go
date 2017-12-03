package main

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/olivere/elastic"
)

type redisClient struct {
	key  string
	pool *redis.Pool
	ec   *elastic.Client
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

func (r *redisClient) read() (string, error) {
	c := r.pool.Get()
	defer c.Close()

	v, err := redis.String(c.Do("LPOP", r.key))
	if err != nil {
		return "", fmt.Errorf("error getting key: %s with: %v", r.key, err)
	}

	return v, nil
}

func (r *redisClient) readFilterAndIndex(msg redis.PMessage) error {
	if string(msg.Data) != r.key {
		return nil
	}
	if !strings.HasSuffix(msg.Channel, "rpush") {
		return nil
	}

	v, err := r.read()
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"original content": v,
	}).Debug("read:")

	filteredJSON, contractName, err := filter(v)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"filtered content": filteredJSON,
	}).Debug("read:")

	err = r.index(fmt.Sprintf("logstash-%s-%d-%d-%d", contractName, time.Now().Year(), time.Now().Month(), time.Now().Day()), filteredJSON)
	if err != nil {
		return err
	}
	return nil
}

func (r *redisClient) consume() error {
	c := r.pool.Get()
	defer c.Close()
	psc := redis.PubSubConn{Conn: c}
	c.Do("CONFIG", "SET", "notify-keyspace-events", "KEA")

	err := psc.PSubscribe("__key*__:*")
	if err != nil {
		return fmt.Errorf("unable to create subscription:%v", err)
	}
	for {
		switch msg := psc.Receive().(type) {
		case redis.Message:
			log.WithFields(log.Fields{"channel": msg.Channel, "data": msg.Data}).Debug("consume: message")
		case redis.PMessage:
			log.WithFields(log.Fields{"channel": msg.Channel, "data": msg.Data, "pattern": msg.Pattern}).Debug("consume: pmessage")
			err = r.readFilterAndIndex(msg)
			if err != nil {
				log.WithFields(log.Fields{"err": err}).Error("consume:")
			}
		case redis.Subscription:
			log.WithFields(log.Fields{"kind": msg.Kind, "channel": msg.Channel, "count": msg.Count}).Debug("consume: subscription")
			if msg.Count == 0 {
				// return
			}
		case error:
			return fmt.Errorf("error: %v", msg)
		}
	}
}
