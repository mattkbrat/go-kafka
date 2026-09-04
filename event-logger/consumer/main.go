package main

import (
	"context"
	"event-consumer/status"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	args := &KafkaArgs{
		Url:       KafkaRoute,
		Topic:     "events",
		Partition: 0,
	}

	s := status.New()
	ProcessMessage(ctx, s, args)
}
