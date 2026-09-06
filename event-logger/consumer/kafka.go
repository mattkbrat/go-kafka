package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/segmentio/kafka-go"
	"google.golang.org/api/option"

	"event-consumer/internal/types"
	"event-consumer/status"
)

const (
	batchSize     = 20
	flushInterval = 3 * time.Second
	shutdownGrace = 5 * time.Second
)

func ProcessMessage(ctx context.Context, s *status.KafkaStatus, config *Config) {
	reader := getKafkaReader(&config.Kafka)
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("error closing kafka reader: %v\n", err)
		}
	}()

	g, _ := json.Marshal(config.Gcp)
	bqClient, err := bigquery.NewClient(ctx, config.Gcp.ProjectId, option.WithAuthCredentialsJSON(option.CredentialsType(config.Gcp.Type), g))
	if err != nil {
		log.Printf("failed to connect to bigquery: %v\n", err)
		s.Set(false)
		return
	}
	defer bqClient.Close()

	log.Printf("connected to bigquery project: %s\n", config.Gcp.ProjectId)
	inserter := bqClient.Dataset(config.DatasetId).Table(config.TableId).Inserter()

	s.Set(true)

	var (
		batch    = make([]types.Event, 0, batchSize)
		toCommit = make([]kafka.Message, 0, batchSize)
		ticker   = time.NewTicker(flushInterval)
	)
	defer ticker.Stop()

	// Flush function commits to BigQuery, then acknowledges Kafka offsets
	flush := func(flushCtx context.Context) error {
		if len(batch) == 0 {
			return nil
		}

		// 1. Insert rows into BigQuery
		if err := inserter.Put(flushCtx, batch); err != nil {
			if putErr, ok := errors.AsType[bigquery.PutMultiError](err); ok {
				for _, rowErr := range putErr {
					log.Printf("bigquery row error (row %d): %v\n", rowErr.RowIndex, rowErr.Errors)
				}
			}
			return fmt.Errorf("bigquery insert failed: %w", err)
		}

		log.Printf("inserted %d rows to bigquery\n", len(batch))

		// 2. Commit offsets to Kafka only after BQ succeeds
		if err := reader.CommitMessages(flushCtx, toCommit...); err != nil {
			return fmt.Errorf("failed to commit kafka offsets: %w", err)
		}

		// 3. Reset buffers
		batch = batch[:0]
		toCommit = toCommit[:0]
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("shutdown signal received, flushing remaining records...")
			s.Set(false)

			// Use a fresh context for final flush and commit
			shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGrace)
			if err := flush(shutdownCtx); err != nil {
				log.Printf("error during shutdown flush: %v\n", err)
			}
			cancel()
			return

		case <-ticker.C:
			if err := flush(ctx); err != nil {
				log.Printf("periodic flush error: %v\n", err)
				s.Set(false)
			} else {
				if !s.IsHealthy() {
					s.Set(true)
				}
			}

		default:
			// Fetch without committing offset automatically
			fetchCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
			m, err := reader.FetchMessage(fetchCtx)
			cancel()

			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					continue
				}
				log.Printf("error fetching message: %v\n", err)
				s.Set(false)
				time.Sleep(1 * time.Second) // backoff on transient errors
				continue
			}

			var event types.Event
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("skipping malformed message at offset %d: %v\n", m.Offset, err)
				log.Print(string(m.Value))
				// Commit skipped poison pills so the queue does not stall
				_ = reader.CommitMessages(ctx, m)
				continue
			}

			batch = append(batch, event)
			toCommit = append(toCommit, m)

			if len(batch) >= batchSize {
				if err := flush(ctx); err != nil {
					log.Printf("batch size flush error: %v\n", err)
					s.Set(false)
				} else {
					s.Set(true)
				}
			}
		}
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
