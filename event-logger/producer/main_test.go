package main

import (
	"bytes"
	"encoding/json/v2"
	"event-logger/handlers"
	"event-logger/internal/data"
	"event-logger/internal/lib"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

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

	u := data.UserParams{
		Name:     "John Doe",
		Username: lib.RandStringRunes(8),
		Email:    "doej@example.com",
		Password: lib.RandStringRunes(12),
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

	// Expected authorization flow
	{
		// Expects created
		{
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if assert.NoError(t, h.Register(c)) {
				assert.Equal(t, http.StatusCreated, rec.Code)
			}
		}

		authorization := ""

		// Can Login
		{
			req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(marshaled))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if err != nil {
				log.Fatalf("Failed to build request: %s", err)
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, h.Login(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)

				r := rec.Result()

				for _, a := range r.Cookies() {
					if a.Name == "authorization" {
						authorization = a.Value
						break
					}
				}
			}

			// gets session
			{
				assert.True(t, len(authorization) > 0)
				req, err := http.NewRequest(http.MethodPost, "/auth/me", bytes.NewReader([]byte{}))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				cookie := new(http.Cookie{
					Name:  "authorization",
					Value: authorization,
				})
				req.AddCookie(cookie)
				if err != nil {
					log.Fatalf("Failed to build request: %s", err)
				}

				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				if assert.NoError(t, h.Me(c)) {
					assert.Equal(t, http.StatusOK, rec.Code)
					assert.NotEmpty(t, rec.Body.String())
					foundUser := data.UserParams{}

					r := rec.Result()
					body, _ := io.ReadAll(r.Body)
					json.Unmarshal(body, &foundUser)

					assert.Equal(t, u.Username, foundUser.Username)
				}

			}

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
