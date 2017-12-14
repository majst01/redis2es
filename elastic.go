package main

import (
	"github.com/google/uuid"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

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
