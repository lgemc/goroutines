package main

import (
	"context"
	"fmt"
	"time"

	patternserver "goroutines-patterns/shared/patterns-server"
)

const (
	numberLength   = 20
	simulatedDelay = 300
)

func genNumbers() []int {
	numbers := make([]int, numberLength)
	for i := 0; i < numberLength; i++ {
		numbers[i] = i
	}

	return numbers
}

func main() {
	numbers := genNumbers()

	fmt.Println("numbers: ", numbers)

	now := time.Now()
	result := make([]int, numberLength)
	// (x^2)
	for i := 0; i < numberLength; i++ {
		response, err := patternserver.PerformOperation(context.Background(), patternserver.OperationPayload{
			Operation: patternserver.Power,
			A:         numbers[i],
			B:         2,
			Delay:     simulatedDelay,
		})
		if err != nil {
			panic(err)
		}

		result[i] = response.Result
	}

	fmt.Println("result: ", result, "elapsed time: ", time.Since(now))
}
