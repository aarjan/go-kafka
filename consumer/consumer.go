package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/aarjan/kafka/model"
	"gopkg.in/olivere/elastic.v5"
)

// NewConsumer returns a new consumer
func NewConsumer(brokerList []string) sarama.Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	master, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		panic(err)
	}
	return master
}

// NewClient returns a new ES client
func NewClient(address, index string) *elastic.Client {
	ctx := context.Background()
	client, err := elastic.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://" + address).Do(ctx)
	if err != nil {
		// Handle error
		log.Fatal(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://" + address)
	if err != nil {
		// Handle error
		log.Fatal(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		// Handle error
		log.Fatal(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(index).BodyString(model.Mapping).Do(ctx)
		if err != nil {
			// Handle error
			log.Fatal(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return client
}
