package executor

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Executor handles sending HTTP requests.
type Executor struct {
	client *http.Client
	debug  bool
}

// NewExecutor creates a new Executor with debug mode option.
func NewExecutor(debug ...bool) *Executor {
	d := false
	if len(debug) > 0 {
		d = debug[0]
	}
	return &Executor{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		debug: d,
	}
}

// Response holds the result of an HTTP request.
type Response struct {
	StatusCode int
	Status     string
	Headers    http.Header
	Body       []byte
	Duration   time.Duration
}

// Do executes the request and returns the response.
func (e *Executor) Do(req *http.Request) (*Response, error) {
	if e.debug {
		fmt.Printf("\n> %s %s\n", req.Method, req.URL.String())
	}

	start := time.Now()
	resp, err := e.client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Body:       body,
		Duration:   duration,
	}, nil
}