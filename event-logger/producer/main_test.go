package main

import (
	"event-logger/handlers"
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
