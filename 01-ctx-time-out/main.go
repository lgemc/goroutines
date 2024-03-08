package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	mathserver "goroutines-patterns/shared/math-server"
)

const (
	numberLength = 20
)

func genNumbers() []int {
	numbers := make([]int, numberLength)
	for i := 0; i < numberLength; i++ {
		numbers[i] = i
	}

	return numbers
}

func main() {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))

	numbers := genNumbers()
	result := make([]int, numberLength)
	fmt.Println("numbers: ", numbers)

	wg := &sync.WaitGroup{}
	for i := 0; i < numberLength; i++ {
		wg.Add(1)
		go func(number int) {
			defer wg.Done()

			operation := mathserver.OperationPayload{
				Operation: mathserver.Power,
				A:         number,
				B:         3,
				Delay:     number * 300,
			}

			response, err := mathserver.PerformOperation(ctx, operation)
			if err != nil {
				fmt.Printf("error: %v\n", err)

				return
			}

			result[number] = response.Result

		}(i)
	}

	wg.Wait()

	fmt.Println("result: ", result)
}
