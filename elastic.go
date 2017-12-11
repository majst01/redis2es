package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

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

// FIXME bulkServices only increase, garbage collect them
// IDEA: keep creation date, remove all older on getBulk, or with separate goroutine
func (r *redisClient) getBulk(indexName string) (*elastic.BulkService, error) {
	var bulk *elastic.BulkService
	if r.indexes[indexName] == nil {
		bulk = r.ec.Bulk().Index(indexName).Type("log")
		r.indexes[indexName] = bulk
		err := r.createIndex(indexName)
		if err != nil {
			return nil, err
		}
	}
	bulk = r.indexes[indexName]
	return bulk, nil
}

func (r *redisClient) getBulks() []*elastic.BulkService {
	var bulks []*elastic.BulkService
	for _, bulk := range r.indexes {
		bulks = append(bulks, bulk)
	}
	return bulks
}

func (r *redisClient) flush() {
	ticker := time.NewTicker(r.bulkTicker)
	for {
		select {
		case <-ticker.C:
			log.Debug("flush: ticker to bulk insert")
			start := time.Now()
			for _, bulk := range r.getBulks() {
				count := bulk.NumberOfActions()
				if count < 1 {
					continue
				}
				res, err := bulk.Do(context.Background())
				if err != nil {
					log.Error(err)
				} else if res.Errors {
					log.Error(fmt.Errorf("bulk commit failed errors:%v", res.Failed()))
				}
				log.Debug("flush: bulk insert res:", res)
				log.WithFields(log.Fields{"duration": time.Now().Sub(start), "count": count}).Info("flush: tick bulk:")
			}
		}
	}
}

func (r *redisClient) index(documents chan document) {
	for {
		select {
		case doc := <-documents:
			start := time.Now()
			log.WithFields(log.Fields{"doc": doc}).Debug("index:")
			bulk, err := r.getBulk(doc.indexName)
			if err != nil {
				log.WithFields(log.Fields{"body": doc.body, "indexName": doc.indexName, "err": err}).Error("index:")
			}
			// FIXME id generation needed ??
			id := base64.URLEncoding.EncodeToString([]byte(doc.body))
			bulk.Add(elastic.NewBulkIndexRequest().Id(id).Doc(doc.body))

			log.WithFields(log.Fields{"outstanding": bulk.NumberOfActions()}).Debug("index:")

			if bulk.NumberOfActions() >= r.bulkSize {
				// Commit
				res, err := bulk.Do(context.Background())
				if err != nil {
					log.WithFields(log.Fields{"error": err}).Error("index: bulk commit")
				} else if res.Errors {
					log.WithFields(log.Fields{"errors": res.Failed()}).Error("index: bulk commit")
				}
				log.WithFields(log.Fields{"result": res}).Debug("index:", res)
				// "bulk" is reset after Do, so you can reuse it
				log.WithFields(log.Fields{"duration": time.Now().Sub(start)}).Info("index: event bulk:")
			}
			log.WithFields(log.Fields{"id": id, "index": doc.indexName}).Debug("index:")
		}
	}
}
