package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	patternserver "goroutines-patterns/shared/patterns-server"
)

const (
	numberLength = 100
	delay        = 500
	maxWorkers   = 5
)

func genNumbers() <-chan int {
	numbers := make(chan int, numberLength)

	go func() {
		defer close(numbers)

		for i := 0; i < numberLength; i++ {
			numbers <- i
		}
	}()

	return numbers
}

func workerPow2(workerId int, in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for number := range in {
		fmt.Printf("processing number %d at worker %d\n", number, workerId)

		operation := patternserver.OperationPayload{
			Operation: patternserver.Power,
			A:         number,
			B:         2,
			Delay:     delay,
		}

		response, err := patternserver.PerformOperation(context.Background(), operation)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		out <- response.Result
	}
}

func main() {
	numbers := genNumbers()
	results := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(maxWorkers)

	now := time.Now()

	for i := 0; i < maxWorkers; i++ {
		go workerPow2(i, numbers, results, &wg)
	}

	resultsArr := make([]int, 0, numberLength)
	go func() {
		for number := range results {
			resultsArr = append(resultsArr, number)
		}
	}()

	wg.Wait()
	close(results)

	fmt.Println("results: ", resultsArr, "elapsed time: ", time.Since(now))
}
