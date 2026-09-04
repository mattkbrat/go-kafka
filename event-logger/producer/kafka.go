package main

import (
	"context"
	"event-logger/internal/status"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func GetKafkaWriter(ctx context.Context, s *status.KafkaStatus, args *KafkaConfig) *kafka.Conn {
	for {

		for i := 1; i < 4; i++ {

			conn, err := kafka.DialLeader(ctx, "tcp", args.Url, args.Topic, args.Partition)
			if err != nil {
				log.Printf("connection error %v", err)
				s.Set(false)

				time.Sleep(time.Second + time.Duration(i*5))
				continue
			}

			return conn
		}

		time.Sleep(time.Second * time.Duration(300))
	}
}
