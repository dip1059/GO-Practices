package main

import (
	"GooglePubSubEx/common"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/satori/go.uuid"
	"google.golang.org/api/option"
	"log"
)


func main() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "bidrace"

	// Creates a client.
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("/home/asha/key.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	topic := common.CreateTopicIfNotExists(client, "BidGenerationRequest")

	id := projectID+"-sub-"+uuid.NewV4().String()
	cfg := pubsub.SubscriptionConfig{
		Topic: topic,
	}

	err = common.CreateSub(client, id, cfg)

	if err != nil {
		log.Println("Failed to Subscribe,",err.Error())
	} else {
		log.Println("Subscribed Successfully")
	}

}
