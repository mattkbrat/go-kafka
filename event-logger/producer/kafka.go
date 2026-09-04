package main

import (
	"context"
	"event-logger/internal/status"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaArgs struct {
	Url       string
	Topic     string
	Partition int
}

func GetKafkaWriter(ctx context.Context, s *status.KafkaStatus, args *KafkaArgs) *kafka.Conn {
	conn, err := kafka.DialLeader(ctx, "tcp", args.Url, args.Topic, args.Partition)
	if err != nil {
		log.Printf("connection error %v", err)
		s.Set(false)
	}

	return conn
}
