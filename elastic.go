package main

import (
	"context"
	"fmt"
)

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
