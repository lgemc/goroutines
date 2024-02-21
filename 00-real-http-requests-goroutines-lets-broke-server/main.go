package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	patternserver "goroutines-patterns/shared/patterns-server"
)

const (
	numberLength   = 200
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

	result := make([]int, numberLength)

	now := time.Now()
	wg := &sync.WaitGroup{}
	// (x^2)
	for i := 0; i < numberLength; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

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
		}(i)
	}

	wg.Wait()

	fmt.Println("result: ", result, "elapsed time: ", time.Since(now))
}
