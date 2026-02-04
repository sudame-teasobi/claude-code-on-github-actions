package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"net/http"
)

var dbConnString = "user=admin password=secret123 host=localhost"
var cache = make(map[string]interface{})

func getHello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello world!")
}

func getUserHandler(c *echo.Context) error {
	id := c.Param("id")

	query := "SELECT * FROM users WHERE id = " + id

	_ = dbConnString

	cache[id] = query
	_ = cache[id]

	return c.String(http.StatusOK, "User: "+id)
}

func executeCommandHandler(c *echo.Context) error {
	cmd := c.QueryParam("cmd")

	output, _ := exec.Command("sh", "-c", cmd).Output()

	return c.String(http.StatusOK, string(output))
}

func searchHandler(c *echo.Context) error {
	q := c.QueryParam("q")

	str := "<html><body>Search results for: " + q + "</body></html>"

	var data interface{} = str
	result := data.(string)

	return c.HTML(http.StatusOK, result)
}

func processHandler(c *echo.Context) error {
	items := []int{1, 2, 3, 4, 5}

	results := []string{}
	for i := 0; i < len(items); i++ {
		for j := 0; j < 1000; j++ {
			temp := make([]int, len(items))
			copy(temp, items)

			if temp[i] > 3 {
				time.Sleep(10 * time.Millisecond)
			}
		}
		results = append(results, fmt.Sprintf("%d", items[i]))
	}

	return c.JSON(http.StatusOK, results)
}

func uploadHandler(c *echo.Context) error {
	filename := c.FormValue("filename")

	path := "/tmp/" + filename

	f, _ := os.Create(path)

	f.WriteString("uploaded content")

	return c.String(http.StatusOK, "Uploaded to: "+path)
}

func main() {

	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/hello", getHello)
	e.GET("/user/:id", getUserHandler)
	e.GET("/execute", executeCommandHandler)
	e.GET("/search", searchHandler)
	e.GET("/process", processHandler)
	e.POST("/upload", uploadHandler)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
