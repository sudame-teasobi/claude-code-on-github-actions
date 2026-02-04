package main

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"net/http"
)

func getHello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello world!")
}

func main() {

	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/hello", getHello)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
