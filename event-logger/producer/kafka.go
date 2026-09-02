package main

import (
	"context"
	"event-logger/internal/status"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

// func connectKafka(ctx context.Context, s *status.KafkaStatus) {
// 	fmt.Println("starting consumer...")
// 	writer := GetKafkaWriter(ctx, s, KafkaRoute, "events")

// 	defer writer.Close()

// 	backoff := time.Second
// 	s.Set(true)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			log.Println("stopping consumer...")
// 			return
// 		default:
// 		}
// 		_, err := writer.WriteMessages(kafka.Message{
// 			Value: []byte("New message!"),
// 		})
// 		if err != nil {
// 			log.Printf("write error %v", err)
// 			s.Set(false)
// 		}
// 	}
// }

func GetKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

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

	// return &kafka.Writer{
	// 	Addr:     kafka.TCP(kafkaURL),
	// 	Topic:    topic,
	// 	Balancer: &kafka.LeastBytes{},
	// }
}
