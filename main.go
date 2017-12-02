package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

type redisClient struct {
	key  string
	pool *redis.Pool
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

	var v string

	reply, err := redis.Values(c.Do("LPOP", r.key))
	if err != nil {
		return "", fmt.Errorf("error getting key: %s with: %v", r.key, err)
	}
	if _, err := redis.Scan(reply, &v); err != nil {
		return "", fmt.Errorf("error scanning result from: %s with: %v", r.key, err)
	}

	return v, nil
}

func (r *redisClient) selectKey(msg redis.PMessage) {
	if string(msg.Data) == r.key {
		v, err := r.read()
		if err != nil {
			fmt.Printf("err:%v", err)
		}
		fmt.Printf("got:%s", v)
	}
}

func (r *redisClient) consume() error {
	c := r.pool.Get()
	defer c.Close()
	psc := redis.PubSubConn{Conn: c}
	c.Do("CONFIG", "SET", "notify-keyspace-events", "KEA")

	psc.PSubscribe("__key*__:*")
	for {
		switch msg := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("Message: %s %s\n", msg.Channel, msg.Data)
		case redis.PMessage:
			fmt.Printf("PMessage: p:%s c:%s d:%s\n", msg.Pattern, msg.Channel, msg.Data)
			r.selectKey(msg)
		case redis.Subscription:
			fmt.Printf("Subscription: %s %s %d\n", msg.Kind, msg.Channel, msg.Count)
			if msg.Count == 0 {
				// return
			}
		case error:
			return fmt.Errorf("error: %v", msg)
		}
	}
}

func main() {
	redisPool := newPool("127.0.0.1", 6379, 0, "", false, false)

	rc := &redisClient{
		pool: redisPool,
		key:  "logstash",
	}

	err := rc.consume()
	if err != nil {
		fmt.Printf("error:%v", err)
	}
}
