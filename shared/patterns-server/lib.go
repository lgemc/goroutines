package patternserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	client    = &http.Client{}
	serverURL = "http://localhost:9999"
)

type Operation string

const (
	Add      Operation = "add"
	Subtract Operation = "subtract"
	Multiply Operation = "multiply"
	Modulus  Operation = "mod"
	Divide   Operation = "divide"
	Power    Operation = "pow"
)

type OperationPayload struct {
	Operation Operation `json:"operation"`
	A         int       `json:"a"`
	B         int       `json:"b"`
	// Delay in milliseconds
	Delay int `json:"delay"`
}

type OperationResult struct {
	Result int `json:"result"`
}

func GetHealth() (*http.Response, error) {
	return client.Get(serverURL + "/health")
}

func PerformOperation(ctx context.Context, operation OperationPayload) (*OperationResult, error) {
	body, err := json.Marshal(operation)
	if err != nil {
		return nil, err
	}

	rawURL := fmt.Sprintf("%s/operate?delay=%d", serverURL, operation.Delay)
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	req := http.Request{
		Method: http.MethodPost,
		URL:    u,
		Body:   io.NopCloser(bytes.NewBuffer(body)),
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}

	response, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var result OperationResult
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return &result, fmt.Errorf("server returned status code %d, response: %s", response.StatusCode, string(b))
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}
