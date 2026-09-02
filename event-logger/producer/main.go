package main

import (
	"context"
	"event-logger/internal/status"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"github.com/labstack/echo/v5"
	"github.com/segmentio/kafka-go"
)

var lock = sync.Mutex{}

type (
	response struct {
		Message string `json:"message,omitempty"`
	}
	handler struct {
		writer *kafka.Conn
	}
)

func (h *handler) healthz(c *echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	if h.writer != nil {
		h.writer.WriteMessages(
			kafka.Message{Value: []byte("healthtest")},
		)
	}

	return c.String(http.StatusOK, "")
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := status.New()

	writer := GetKafkaWriter(ctx, s, &KafkaArgs{
		Url:       KafkaRoute,
		Topic:     "events",
		Partition: 0,
	})

	e := echo.New()
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Hello, World!"})
	})

	h := handler{
		writer,
	}

	e.GET("/healthz", h.healthz)
	e.Start(":1323")
}
