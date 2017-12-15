package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

// ElasticClient is used to store data in elasticsearch
type ElasticClient struct {
	client        *elastic.Client
	bulkProcessor *elastic.BulkProcessor
}

// NewElasticClient create a new instance of a elasticClient
func NewElasticClient(spec Specification) *ElasticClient {
	client, err := elastic.NewSimpleClient(elastic.SetURL(spec.ElasticURLs...))
	if err != nil {
		log.WithFields(log.Fields{"error connecting to elastic": err}).Error("main:")
	}
	bulk, err := client.BulkProcessor().
		Name("BackgroundWorker-1").
		Workers(spec.PoolSize).         // number of workers
		BulkActions(spec.BulkSize).     // commit if # requests >= BulkSize
		BulkSize(2 << 20).              // commit if size of requests >= 2 MB
		FlushInterval(spec.BulkTicker). // commit every given interval
		Stats(true).                    // collect stats
		Do(context.Background())
	if err != nil {
		log.WithFields(log.Fields{"error creating bulkprocessor": err}).Fatal("main:")
	}

	ec := &ElasticClient{
		client:        client,
		bulkProcessor: bulk,
	}
	return ec
}

func (e *ElasticClient) close() {
	e.bulkProcessor.Close()
	e.client.Stop()
}

func (e *ElasticClient) index(documents chan document) {
	for {
		select {
		case doc := <-documents:
			log.WithFields(log.Fields{"doc": doc}).Debug("index:")

			id := uuid.New().String()
			request := elastic.NewBulkIndexRequest().Index(doc.indexName).Type("log").Id(id).Doc(doc.body)
			e.bulkProcessor.Add(request)

			log.WithFields(log.Fields{"id": id, "index": doc.indexName}).Debug("index:")
		}
	}
}

func (e *ElasticClient) stats() {
	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-ticker.C:
			stats := e.bulkProcessor.Stats()

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
