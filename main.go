package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/json-iterator/go"
	"github.com/olivere/elastic"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
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

func (r *redisClient) filter(input string) (string, string, error) {
	contractName := "catchall"
	data := make(map[string]interface{})
	err := json.UnmarshalFromString(input, &data)
	if err != nil {
		return "", contractName, fmt.Errorf("cannot decode data:%v", err)
	}

	// check if contract in any case is present, lowercase then
	for k, v := range data {
		if strings.ToLower(k) == "contract" {
			oldValue := strings.ToLower(v.(string))
			delete(data, k)
			data["contract"] = oldValue
			contractName = oldValue
		}
	}

	result, err := json.MarshalToString(&data)
	if err != nil {
		return "", contractName, fmt.Errorf("cannot encode data:%v", err)
	}
	return result, contractName, nil

}

func (r *redisClient) index(name, bodyJSON string) error {
	exists, err := r.ec.IndexExists(name).Do(context.Background())
	if err != nil {
		return fmt.Errorf("index:%s cannot be checked:%v", name, err)
	}
	if !exists {
		createIndex, err := r.ec.CreateIndex(name).Do(context.Background())
		if err != nil {
			return fmt.Errorf("cannot create index:%s %v", name, err)
		}
		if !createIndex.Acknowledged {
			return fmt.Errorf("create index:%s was not acknowledged", name)
		}
		fmt.Printf("index:%s created\n", name)
	}

	writeIndex, err := r.ec.Index().Index(name).Type("log").BodyJson(bodyJSON).Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot add %s to index %s err:%v", bodyJSON, name, err)
	}
	fmt.Printf("indexed: id:%s index:%s type:%s\n", writeIndex.Id, writeIndex.Index, writeIndex.Type)
	return nil
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

func (r *redisClient) selectKey(msg redis.PMessage) {
	if string(msg.Data) == r.key && strings.HasSuffix(string(msg.Channel), "rpush") {
		v, err := r.read()
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
		fmt.Printf("got:%s\n", v)

		filteredJSON, contractName, err := r.filter(v)
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
		err = r.index(fmt.Sprintf("logstash-%s-%d-%d-%d", contractName, time.Now().Year(), time.Now().Month(), time.Now().Day()), filteredJSON)
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
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
		fmt.Printf("error:%v", err)
	}

}
