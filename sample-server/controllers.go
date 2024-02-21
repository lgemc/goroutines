package main

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

type binaryOperation struct {
	Operation string `json:"operation"`
	A         int    `json:"a"`
	B         int    `json:"b"`
}

func operate(c echo.Context) error {
	var operation binaryOperation
	err := c.Bind(&operation)
	if err != nil {
		return c.String(http.StatusBadRequest, "Error binding body")
	}

	// Simulate delay
	delay := c.QueryParam("delay")
	if delay != "" {
		intDelay, err := strconv.Atoi(delay)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid delay")
		}

		if intDelay > 0 {
			<-time.After(time.Duration(intDelay) * time.Millisecond)
		}
	}

	var result int
	switch operation.Operation {
	case "add":
		result = operation.A + operation.B
	case "subtract":
		result = operation.A - operation.B
	case "multiply":
		result = operation.A * operation.B
	case "divide":
		result = operation.A / operation.B
	case "pow":
		result = int(math.Pow(float64(operation.A), float64(operation.B)))
	case "mod":
		result = operation.A % operation.B
	default:
		return c.String(http.StatusBadRequest, "Invalid operation")
	}

	return c.JSON(http.StatusOK, map[string]int{"result": result})
}
