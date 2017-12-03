package main

import (
	"fmt"

	"github.com/olivere/elastic"
)

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
