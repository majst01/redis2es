package main

import (
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
			log.WithFields(
				log.Fields{
					"flushed":   stats.Flushed,
					"commited":  stats.Committed,
					"indexed":   stats.Indexed,
					"created":   stats.Created,
					"updated":   stats.Updated,
					"succeeded": stats.Succeeded,
					"failed":    stats.Failed,
				}).Info("stats:")

			//	fmt.Printf("Number of times flush has been invoked: %d\n", stats.Flushed)
			//	fmt.Printf("Number of times workers committed reqs: %d\n", stats.Committed)
			//	fmt.Printf("Number of requests indexed            : %d\n", stats.Indexed)
			//	fmt.Printf("Number of requests reported as created: %d\n", stats.Created)
			//	fmt.Printf("Number of requests reported as updated: %d\n", stats.Updated)
			//	fmt.Printf("Number of requests reported as success: %d\n", stats.Succeeded)
			//	fmt.Printf("Number of requests reported as failed : %d\n", stats.Failed)

			for i, w := range stats.Workers {
				//		fmt.Printf("Worker %d: Number of requests queued: %d\n", i, w.Queued)
				//		fmt.Printf("           Last response time       : %v\n", w.LastDuration)

				// FIXME put these into above log line as nested struct
				log.WithFields(
					log.Fields{
						"worker":       i,
						"queued":       w.Queued,
						"lastduration": w.LastDuration,
					}).Info("stats:")

			}
		}
	}
}
