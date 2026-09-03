package handlers

import (
	"errors"
	"event-logger/internal/data"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/segmentio/kafka-go"
)

var (
	ErrMissingAttr      = errors.New("Username, email, and password are required")
	ErrInsecurePassword = errors.New("Password must be at least 8 characters")
	ErrUserTaken        = errors.New("Username and email must be unique")
)

func AddUser(user *data.UserType) error {

	if user.Email == "" || user.Username == "" || user.Password == "" {
		return ErrMissingAttr
	}

	if len(user.Password) < 8 {
		return ErrInsecurePassword
	}

	for _, u := range data.Users {
		if u.Username == user.Username || u.Email == user.Email {
			return ErrUserTaken
		}
	}

	data.Users[data.Seq] = user

	return nil
}

func (h *Handler) Register(c *echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	var user data.UserType

	hasWriter := h.Writer != nil

	if err := c.Bind(&user); err != nil {
		if hasWriter {
			h.Writer.WriteMessages(kafka.Message{
				Value: fmt.Appendf(nil, "Failed to register user %v", err.Error()),
			})
		}
		return c.String(http.StatusBadRequest, "")
	}

	if err := AddUser(&user); err != nil {
		m := err.Error()
		if hasWriter {
			h.Writer.WriteMessages(kafka.Message{
				Value: fmt.Appendf(nil, "Failed to register user %v", m),
			})
		}
		return c.String(http.StatusUnprocessableEntity, m)
	}

	if hasWriter {
		h.Writer.WriteMessages(kafka.Message{
			Value: fmt.Appendf(nil, "registered %s", user.Email),
		})
	}

	return c.String(http.StatusCreated, "")
}
