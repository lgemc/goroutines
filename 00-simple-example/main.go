package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"goroutines-patterns/shared/patterns-server"
)

const (
	maxWorkers = 10
)

func getHealth(workerId int) (string, error) {
	response, err := patternserver.GetHealth()
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {

		return "", fmt.Errorf("server is not healthy, status: %d, body: %s", response.StatusCode, string(body))
	}

	return fmt.Sprintf("response collected at worker %d, response: %s", workerId, string(body)), nil
}

func main() {
	wg := &sync.WaitGroup{}
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			response, err := getHealth(workerId)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(response)
		}(i)
	}

	wg.Wait()
}
