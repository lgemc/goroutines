package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	mathserver "goroutines-patterns/shared/math-server"
)

const (
	numberLength   = 20
	simulatedDelay = 500
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

			response, err := mathserver.PerformOperation(context.Background(), mathserver.OperationPayload{
				Operation: mathserver.Power,
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
