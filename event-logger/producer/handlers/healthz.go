package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/segmentio/kafka-go"
)

func (h *Handler) Healthz(c *echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	if h.Writer != nil {
		h.Writer.WriteMessages(
			kafka.Message{Value: []byte("healthtest")},
		)
	}

	return c.String(http.StatusOK, "")
}
