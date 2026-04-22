package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gongdaowen/apix/internal/executor"
	"github.com/tidwall/pretty"
)

// OutputMode defines the response output format.
type OutputMode int

const (
	OutputPretty   OutputMode = iota // Default: human-readable with colors
	OutputJSON                       // Full response as JSON
	OutputRaw                        // Only raw body
)

// ResponseJSON represents the complete response in JSON format.
type ResponseJSON struct {
	StatusCode int               `json:"status_code"`
	Status     string            `json:"status"`
	Duration   string            `json:"duration_ms"`
	URL        string            `json:"url,omitempty"`
	Method     string            `json:"method,omitempty"`
	Headers    map[string]string `json:"headers"`
	Body       json.RawMessage   `json:"body"`
}

// Formatter handles formatting of HTTP responses.
type Formatter struct {
	colorEnabled bool
	mode         OutputMode
}

func NewFormatter() *Formatter {
	return &Formatter{
		colorEnabled: true,
		mode:         OutputPretty,
	}
}

// SetMode sets the output mode.
func (f *Formatter) SetMode(mode OutputMode) {
	f.mode = mode
}

// Print outputs the response based on the configured mode.
func (f *Formatter) Print(resp *executor.Response) error {
	switch f.mode {
	case OutputRaw:
		return f.PrintRaw(resp)
	case OutputJSON:
		return f.PrintJSON(resp)
	default:
		return f.PrintPretty(resp)
	}
}

// PrintPretty outputs the response in a human-readable format with colors.
func (f *Formatter) PrintPretty(resp *executor.Response) error {
	// Status line
	statusColor := f.statusColor(resp.StatusCode)
	fmt.Fprintf(os.Stdout, "\n%s %s\n", statusColor, resp.Status)
	fmt.Fprintf(os.Stdout, "Duration: %s\n", resp.Duration.Round(time.Millisecond))

	// Headers
	fmt.Fprintln(os.Stdout, "\n─── Response Headers ───")
	for key, values := range resp.Headers {
		for _, value := range values {
			fmt.Fprintf(os.Stdout, "  %s: %s\n", key, value)
		}
	}

	// Body
	if len(resp.Body) > 0 {
		fmt.Fprintln(os.Stdout, "\n─── Response Body ───")
		prettyBody := pretty.Pretty(resp.Body)
		fmt.Fprintf(os.Stdout, "%s\n", prettyBody)
	}

	fmt.Fprintln(os.Stdout)

	// Exit with error code for non-2xx responses
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request returned %s", resp.Status)
	}

	return nil
}

// PrintRaw outputs only the raw response body (no headers, no extra text).
func (f *Formatter) PrintRaw(resp *executor.Response) error {
	if len(resp.Body) > 0 {
		// Try to pretty-print JSON if it's valid JSON
		var jsonData interface{}
		if err := json.Unmarshal(resp.Body, &jsonData); err == nil {
			// Valid JSON, pretty print it
			prettyBody, _ := json.MarshalIndent(jsonData, "", "  ")
			os.Stdout.Write(prettyBody)
			os.Stdout.WriteString("\n")
		} else {
			// Not JSON, print as-is
			os.Stdout.Write(resp.Body)
		}
	}

	// Exit with error code for non-2xx responses
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request returned %s", resp.Status)
	}

	return nil
}

// PrintJSON outputs the full response (status, headers, body) as a single JSON object.
func (f *Formatter) PrintJSON(resp *executor.Response) error {
	headers := make(map[string]string)
	for key, values := range resp.Headers {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	body := resp.Body
	// Ensure body is valid JSON for the JSON output
	if len(body) == 0 {
		body = []byte("null")
	} else {
		// Try to validate/normalize the JSON
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err == nil {
			body, _ = json.Marshal(jsonData)
		}
	}

	response := ResponseJSON{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Duration:   fmt.Sprintf("%.3f", resp.Duration.Seconds()*1000),
		Headers:    headers,
		Body:       json.RawMessage(body),
	}

	output, _ := json.MarshalIndent(response, "", "  ")
	fmt.Fprintln(os.Stdout, string(output))

	// Exit with error code for non-2xx responses
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request returned %s", resp.Status)
	}

	return nil
}

// statusColor returns colored status code string.
func (f *Formatter) statusColor(code int) string {
	if !f.colorEnabled {
		return fmt.Sprintf("%d", code)
	}
	
	switch {
	case code >= 200 && code < 300:
		return fmt.Sprintf("\033[32m%d\033[0m", code) // Green
	case code >= 300 && code < 400:
		return fmt.Sprintf("\033[33m%d\033[0m", code) // Yellow
	case code >= 400 && code < 500:
		return fmt.Sprintf("\033[31m%d\033[0m", code) // Red
	case code >= 500:
		return fmt.Sprintf("\033[35m%d\033[0m", code) // Purple
	default:
		return fmt.Sprintf("%d", code)
	}
}

// PrintCurl generates a curl command equivalent for the request.
func (f *Formatter) PrintCurl(method, url string, headers map[string]string, body []byte) {
	fmt.Fprintf(os.Stdout, "\ncurl -X %s \\\n", method)
	fmt.Fprintf(os.Stdout, "  '%s'", url)

	for key, value := range headers {
		fmt.Fprintf(os.Stdout, " \\\n  -H '%s: %s'", key, value)
	}

	if len(body) > 0 {
		fmt.Fprintf(os.Stdout, " \\\n  -d '%s'", string(body))
	}

	fmt.Fprintln(os.Stdout)
}
