package main

import (
	"GooglePubSubEx/common"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	//"github.com/satori/go.uuid"
	"google.golang.org/api/option"
	"log"
)

type Data struct {
	UserID string	`json:"user_id"`
	BidsAmount string	`json:"bids_amount"`
}

func main() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "bidrace"

	// Creates a client.
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("./bidrace-fa1799950676.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	topic := common.CreateTopicIfNotExists(client, "BidGenerationRequest")

	id := "Bidstore"
	cfg := pubsub.SubscriptionConfig{
		Topic: topic,
	}

	err = common.CreateSub(client, id, cfg)

	if err != nil {
		log.Println("Failed to Subscribe,",err.Error())
	} else {
		log.Println("Subscribed Successfully")
	}

	data, _ := json.Marshal(map[string]string{
		"user_id":"3",
		"bids_amount":"2000",
	})
	log.Println(string(data))
	pubResult := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	log.Println(pubResult.Get(ctx))
	sub := client.Subscription("Bidstore")
	err = sub.Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
		message.Ack()
		var d Data
		err = json.Unmarshal(message.Data, &d)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println(string(message.Data))
		log.Println("Got message:", d.UserID, d.BidsAmount)
	})
	if err != nil {
		log.Println(err.Error())
	}
}
