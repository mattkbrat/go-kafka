package handlers

import (
	"sync"

	"github.com/segmentio/kafka-go"
)

type Handler struct {
	Writer *kafka.Conn
}

var lock = sync.Mutex{}
