package main

import (
	"bytes"
	"encoding/json/v2"
	"event-logger/handlers"
	"event-logger/internal/data"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

// Source - https://stackoverflow.com/a/31832326
// Posted by icza, modified by community. See post 'Timeline' for change history
// Retrieved 2026-09-03, License - CC BY-SA 4.0

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Test_healthz(t *testing.T) {
	e := echo.New()
	req, _ := http.NewRequest(http.MethodGet, "/healthz", strings.NewReader(""))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &handlers.Handler{
		Writer: nil,
	}

	if assert.NoError(t, h.Healthz(c)) {

		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func Test_Register(t *testing.T) {
	e := echo.New()

	endpoint := "/auth/register"

	h := &handlers.Handler{
		Writer: nil,
	}

	{
		// Expects failure from no body
		req, _ := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(""))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, h.Register(c)) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		}
	}

	u := data.UserType{
		Name:     "John Doe",
		Username: RandStringRunes(8),
		Email:    "doej@example.com",
		Password: RandStringRunes(12),
	}

	marshaled, err := json.Marshal(u)

	if err != nil {
		log.Fatalf("Failed to marshal user: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(marshaled))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err != nil {
		log.Fatalf("Failed to build request: %s", err)
	}

	{
		// Expects created
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, h.Register(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	}

	{
		// Expects bad request from duplicate user params
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, h.Register(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	}

}
