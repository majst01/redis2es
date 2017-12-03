package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
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
		log.WithFields(log.Fields{"created index:": name}).Debug("index:")
	}

	writeIndex, err := r.ec.Index().Index(name).Type("log").BodyJson(bodyJSON).Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot add %s to index %s err:%v", bodyJSON, name, err)
	}
	log.WithFields(log.Fields{"id": writeIndex.Id, "index": writeIndex.Index, "type": writeIndex.Type}).Debug("index:")

	return nil
}
