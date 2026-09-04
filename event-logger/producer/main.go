package main

import (
	"context"
	"event-logger/handlers"
	"event-logger/internal/status"
	"log"
	"os"
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

	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		log.Fatal("Must provide CONFIG_FILE env variable")
	}

	config := ReadConfig(configFile)

	writer := GetKafkaWriter(ctx, s, &config.Kafka)

	e := echo.New()
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Hello, World!"})
	})

	h := handlers.Handler{
		Writer: writer,
	}

	e.GET("/healthz", h.Healthz)
	auth := e.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/me", h.Me)
	e.Start(":1323")
}
