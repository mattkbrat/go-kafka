package status

import (
	"log"
	"sync/atomic"
)

// KafkaStatus is a thread-safe health flag, shared between the
// consumer goroutine and the HTTP handlers.
type KafkaStatus struct {
	healthy atomic.Bool
}

func New() *KafkaStatus { return &KafkaStatus{} }

func (s *KafkaStatus) Set(ok bool) {
	log.Printf("kafka: %v", ok)
	s.healthy.Store(ok)
}
func (s *KafkaStatus) IsHealthy() bool { return s.healthy.Load() }
