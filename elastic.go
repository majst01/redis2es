package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

func (r *redisClient) createIndex(name string) error {
	exists, err := r.ec.IndexExists(name).Do(context.Background())
	if err != nil {
		return fmt.Errorf("index:%s cannot be checked:%v", name, err)
	}
	if exists {
		return nil
	}
	createIndex, err := r.ec.CreateIndex(name).Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot create index:%s %v", name, err)
	}
	if !createIndex.Acknowledged {
		return fmt.Errorf("create index:%s was not acknowledged", name)
	}
	log.WithFields(log.Fields{"created index:": name}).Debug("index:")
	return nil
}

func (r *redisClient) index(documents chan document) {
	for {
		select {
		case doc := <-documents:
			log.WithFields(log.Fields{"doc": doc}).Debug("index:")

			id := uuid.New().String()
			request := elastic.NewBulkIndexRequest().Index(doc.indexName).Type("log").Id(id).Doc(doc.body)
			r.bulkProcessor.Add(request)

			log.WithFields(log.Fields{"id": id, "index": doc.indexName}).Debug("index:")
		}
	}
}
