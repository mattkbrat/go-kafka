package main

import (
	"context"
	"event-logger/handlers"
	"event-logger/internal/status"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v5"
)

type (
	response struct {
		Message string `json:"message,omitempty"`
	}
)

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

	h := handlers.Handler{
		Writer: writer,
	}

	e.GET("/healthz", h.Healthz)
	health := e.Group("/auth")
	health.POST("/register", h.Register)
	e.Start(":1323")
}
