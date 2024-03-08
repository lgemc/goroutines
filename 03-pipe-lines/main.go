package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	mathserver "goroutines-patterns/shared/math-server"
)

const (
	numberLength = 100
	delay        = 500
	maxWorkers   = 20

	// with 100 numbers, 500ms delay and 1 worker, the elapsed time are 50s per operation
	// with 100 numbers, 500ms delay and 5 workers, the elapsed time are 10s
	// with 100 numbers, 500ms delay and 10 workers, the elapsed time are 5s

	// suppose we have 3 operations like , 4(x^2) + 2
	// without pipelines, with 100 numbers, 500ms delay and 1 worker
	// we will have 300 operations each one 500ms, the elapsed time will be 150s (2.5 minutes)

	// pipelines
	// if we have three operations like , 4(x^2) + 2: with 100 numbers, 500ms delay and 10 workers, the elapsed time
	// is approximately ~5s (between 5 and 6 seconds)
	// it will take 5s also
	//
	// With one worker, the elapsed time is 50s per operation, without pipelines with 3 operations like
	// example above, it will be 150s (2.5 minutes)
)

func genNumbers() <-chan pipeMessage {
	numbers := make(chan pipeMessage, numberLength)

	go func() {
		defer close(numbers)

		for i := 0; i < numberLength; i++ {
			numbers <- pipeMessage{
				xOriginal: i,
				x:         i,
			}
		}
	}()

	return numbers
}

type pipeMessage struct {
	xOriginal int
	x         int
	y         int
	err       error
}

func invertXY(in <-chan pipeMessage) <-chan pipeMessage {
	out := make(chan pipeMessage, numberLength)

	wg := &sync.WaitGroup{}

	multiplex := func() {
		defer wg.Done()

		for message := range in {
			out <- pipeMessage{
				xOriginal: message.xOriginal,
				x:         message.y,
				y:         message.x,
				err:       message.err,
			}
		}
	}

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go multiplex()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func operate(operation mathserver.Operation, operand int, in <-chan pipeMessage) <-chan pipeMessage {
	out := make(chan pipeMessage, numberLength)

	wg := &sync.WaitGroup{}
	multiplex := func(workerId int) {
		defer wg.Done()

		for message := range in {
			fmt.Printf("processing number %d at worker %d\n", message.x, workerId)

			operation := mathserver.OperationPayload{
				Operation: operation,
				A:         message.x,
				B:         operand,
				Delay:     delay,
			}

			response, err := mathserver.PerformOperation(context.Background(), operation)
			if err != nil {
				out <- pipeMessage{
					xOriginal: message.xOriginal,
					x:         message.x,
					err:       err,
				}

				continue
			}

			out <- pipeMessage{
				xOriginal: message.xOriginal,
				x:         message.x,
				y:         response.Result,
				err:       err,
			}
		}
	}

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)

		go multiplex(i)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func printResults(resultsArr []pipeMessage) {
	for _, message := range resultsArr {
		if message.err != nil {
			fmt.Printf("error at x %d: %v\n", message.xOriginal, message.err)

			continue
		}

		fmt.Printf("result, x: %d, y: %d\n", message.xOriginal, message.y)
	}
}

func main() {
	numbers := genNumbers()
	now := time.Now()

	// 4(x^2) + 2

	// Dividimos la operacion en 3 partes
	// (x^2)
	result := operate(mathserver.Power, 2, numbers)
	result = invertXY(result)

	// 4(x^2)
	result = operate(mathserver.Multiply, 4, result)
	result = invertXY(result)

	// 4(x^2) + 2
	result = operate(mathserver.Add, 2, result)

	resultsArr := make([]pipeMessage, 0, numberLength)
	for message := range result {
		if message.err != nil {
			fmt.Printf("error: %v\n", message.err)

			continue
		}

		resultsArr = append(resultsArr, message)
	}

	printResults(resultsArr)
	plotMessages(resultsArr)

	fmt.Printf("elapsed time: %v\n", time.Since(now))
}
