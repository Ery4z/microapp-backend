package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger()) // This will log the incoming requests

	e.Static("/", "static")
	// Routes
	e.POST("/data", postData)
	e.POST("/groups", createGroup)
	e.PUT("/groups/:id", updateGroup)
	e.GET("/groups/:id", getGroup)
	e.GET("/groups", listGroups)
	e.GET("/groups/:groupId/sensors", getSensors)

	// Database initialization
	initDB()

	e.Logger.Fatal(e.Start(":1323"))
}
