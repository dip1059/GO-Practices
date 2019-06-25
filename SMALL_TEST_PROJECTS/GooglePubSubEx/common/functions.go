// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command topics is a tool to manage Google Cloud Pub/Sub topics by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.

// This file is modified by Ta-Ching Chen (https://tachingchen.com/) on 2017.

package common

import (
	"fmt"
	"google.golang.org/api/option"
	"log"

	"cloud.google.com/go/pubsub"
	//"cloud.google.com/go/storage"
	//"golang.org/x/net/context"
	"context"
)

func CreateClient(projectID string) *pubsub.Client {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID, option.WithServiceAccountFile("/home/asha/Bidrace-cf6ecfd807ad.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func CreateTopicIfNotExists(client *pubsub.Client, name string) *pubsub.Topic {
	ctx := context.Background()
	topic := client.Topic(name)
	isExist, err := topic.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if isExist {
		return topic
	}
	topic, err = client.CreateTopic(ctx, name)
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}
	return topic
}

func CreateSub(client *pubsub.Client, id string, cfg pubsub.SubscriptionConfig) error {
	ctx := context.Background()
	// [START create_subscription]
	sub, err := client.CreateSubscription(ctx, id, cfg)
	if err != nil {
		return err
	}
	fmt.Printf("Created subscription: %v\n", sub)
	// [END create_subscription]
	return nil
}

func DeleteSub(client *pubsub.Client, name string) error {
	ctx := context.Background()
	// [START delete_subscription]
	sub := client.Subscription(name)
	if err := sub.Delete(ctx); err != nil {
		return err
	}
	fmt.Println("Subscription deleted.")
	// [END delete_subscription]
	return nil
}
