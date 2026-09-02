package main

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v5"
)

var lock = sync.Mutex{}

type (
	response struct {
		Message string `json:"message,omitempty"`
	}
	handler struct {
	}
)

func (h *handler) healthz(c *echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	// r := response{
	// Message: "Ok",
	// }

	// if r.Message != "" {
	// 	return c.JSON(http.StatusOK, r)
	// }

	return c.String(http.StatusOK, "")

}

func main() {
	e := echo.New()
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Hello, World!"})
	})

	h := handler{}

	e.GET("/healthz", h.healthz)
	e.Start(":1323")
}
