package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	mathserver "goroutines-patterns/shared/math-server"
)

const (
	// con 100 numeros, 1 worker (sin gorutinas) y delay de 500, se demoraria 50 segundos
	// con 100 numeros, 5 workers (con gorutinas) y delay de 500, se demoraria 10 segundos
	// con 100 numeros, 10 workers (con gorutinas) y delay de 500, se demoraria 5 segundos
	numberLength = 100
	delay        = 500
	maxWorkers   = 20
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

		operation := mathserver.OperationPayload{
			Operation: mathserver.Power,
			A:         number,
			B:         2,
			Delay:     delay,
		}

		response, err := mathserver.PerformOperation(context.Background(), operation)
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
