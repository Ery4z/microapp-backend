package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Routes
	e.GET("/", helloWorld)
	e.POST("/data", postData)
	e.POST("/groups", createGroup)
	e.PUT("/groups/:id", updateGroup)
	e.GET("/groups/:id", getGroup)
	e.GET("/groups", listGroups)

	// Database initialization
	initDB()

	e.Logger.Fatal(e.Start(":1323"))
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
