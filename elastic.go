package main

import (
	"fmt"
	"time"

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

func (r *redisClient) stats() {
	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-ticker.C:
			stats := r.bulkProcessor.Stats()

			//	fmt.Printf("Number of times flush has been invoked: %d\n", stats.Flushed)
			//	fmt.Printf("Number of times workers committed reqs: %d\n", stats.Committed)
			//	fmt.Printf("Number of requests indexed            : %d\n", stats.Indexed)
			//	fmt.Printf("Number of requests reported as created: %d\n", stats.Created)
			//	fmt.Printf("Number of requests reported as updated: %d\n", stats.Updated)
			//	fmt.Printf("Number of requests reported as success: %d\n", stats.Succeeded)
			//	fmt.Printf("Number of requests reported as failed : %d\n", stats.Failed)

			fields := log.Fields{
				"flushed":   stats.Flushed,
				"committed": stats.Committed,
				"indexed":   stats.Indexed,
				"created":   stats.Created,
				"updated":   stats.Updated,
				"succeeded": stats.Succeeded,
				"failed":    stats.Failed,
			}

			for i, w := range stats.Workers {
				//		fmt.Printf("Worker %d: Number of requests queued: %d\n", i, w.Queued)
				//		fmt.Printf("           Last response time       : %v\n", w.LastDuration)
				fields[fmt.Sprintf("w%d.queued", i)] = w.Queued
				fields[fmt.Sprintf("w%d.lastduration", i)] = w.LastDuration
			}

			log.WithFields(fields).Info("stats:")
		}
	}
}
