package main

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"goroutines-patterns/shared/env"
)

var (
	port = fmt.Sprintf(":%s", env.GetString("PORT", "9999"))
)

func newServer() *echo.Echo {
	e := echo.New()
	e.GET("/health", health)
	e.POST("/operate", operate)

	return e
}

func start() {
	e := newServer()

	err := e.Start(port)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Info("Server started on port " + port)
}
