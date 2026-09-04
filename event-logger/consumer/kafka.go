package main

import (
	"context"
	"errors"
	"event-consumer/status"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func timeIn(name string) time.Time {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return time.Now().In(loc)
}

func ProcessMessage(
	ctx context.Context,
	s *status.KafkaStatus,
	args *KafkaConfig,
) {

	conn := getKafkaReader(args)

	if conn == nil {
		s.Set(false)
		return
	}

	loc, _ := time.LoadLocation("America/Denver")
	backoff := time.Second

	for {
		select {
		case <-ctx.Done():
			log.Println("stopping consumer...")
			if err := conn.Close(); err != nil {
				log.Println("failed to close connection:", err)
				s.Set(false)
			}
			return
		default:
		}
		if !s.IsHealthy() {
			log.Println("setting healthy")
			s.Set(true)
		}

		// batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

		m, err := conn.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("read error: %v", err)
			s.Set(false)
			if backoff < 30*time.Second {
				backoff *= 20
			}
			time.Sleep(backoff)
			continue
		}

		log.Printf("%s\n %s\n\n", m.Time.In(loc), string(m.Value))
	}
}

func getKafkaReader(args *KafkaConfig) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{args.Url},
		GroupID:  args.GroupId,
		Topic:    args.Topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}
