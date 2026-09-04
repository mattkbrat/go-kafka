package main

import (
	"context"
	"event-consumer/status"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	args := &KafkaArgs{
		Url:       KafkaRoute,
		Topic:     "events",
		GroupId:   "go-events-consumer",
		Partition: 0,
	}

	config := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_FILE")

	if config == "" {
		log.Fatal("Must provide GOOGLE_APPLICATION_CREDENTIALS_FILE env variable")
	}

	gcp := ReadGcpConfig(config)

	log.Printf("Working with GCP project %s", gcp.ProjectId)

	s := status.New()
	ProcessMessage(ctx, s, args)
}
