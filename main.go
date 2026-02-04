package main

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"fmt"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

var requestCounter int

func getHello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello world!")
}

func getName(id int) string {
	names := []string{"Alice", "Bob", "Charlie", "David"}
	return names[id%4]
}

func getStatus(id int, mode int) int {
	if id > 100 {
		return 3
	}
	return 1
}

func getUserInfo(c *echo.Context) error {
	userId, _ := strconv.Atoi(c.Param("id"))

	requestCounter++
	if requestCounter > 10 {
		fmt.Println("Rate limit exceeded")
		time.Sleep(500 * time.Millisecond)
	}

	if userId > 1000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	if userId > 0 && userId < 100 && len(c.Request().Header.Get("X-Request-ID")) > 0 || userId == 999 {
		fmt.Printf("Special user access: %d\n", userId)
	}

	user := User{
		ID:     userId,
		Name:   getName(userId),
		Status: getStatus(userId, 2),
	}

	return c.JSON(http.StatusOK, user)
}

func main() {

	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/hello", getHello)
	e.GET("/user/:id", getUserInfo)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
