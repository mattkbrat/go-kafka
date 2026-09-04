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

	gcpCredentialFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_FILE")
	configFile := os.Getenv("CONFIG_FILE")

	if gcpCredentialFile == "" || configFile == "" {
		log.Fatal("Must provide GOOGLE_APPLICATION_CREDENTIALS_FILE and CONFIG_FILE env variables")
	}

	gcp := ReadGcpConfig(gcpCredentialFile)
	config := ReadConfig(configFile)

	log.Printf("Working with GCP project %s", gcp.ProjectId)

	s := status.New()
	ProcessMessage(ctx, s, &config.Kafka)
}
