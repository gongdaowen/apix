package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/gongdaowen/apix/internal/executor"
)

// OutputMode defines the response output format.
type OutputMode int

const (
	OutputEnhanced  OutputMode = iota // Default: enhanced body only
	OutputFull                        // Full response with enhanced body
	OutputRaw                         // Body as pretty-printed JSON
	OutputBodyJSON                    // Body as JSON (compact)
	OutputFullJSON                    // Full response as JSON
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
		mode:         OutputEnhanced, // Default to enhanced body only
	}
}

// SetMode sets the output mode.
func (f *Formatter) SetMode(mode OutputMode) {
	f.mode = mode
}

// Print outputs the response based on the configured mode.
func (f *Formatter) Print(resp *executor.Response) error {
	switch f.mode {
	case OutputFull:
		return f.PrintFull(resp)
	case OutputRaw:
		return f.PrintRaw(resp)
	case OutputBodyJSON:
		return f.PrintBodyJSON(resp)
	case OutputFullJSON:
		return f.PrintFullJSON(resp)
	default:
		// Default: OutputEnhanced - show only body with enhanced formatting
		return f.PrintEnhanced(resp)
	}
}

// PrintFull outputs the complete response with enhanced body formatting.
func (f *Formatter) PrintFull(resp *executor.Response) error {
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

	// Body - Enhanced display
	if len(resp.Body) > 0 {
		fmt.Fprintln(os.Stdout, "\n─── Response Body ───")
		f.printEnhancedBody(resp.Body)
	}

	fmt.Fprintln(os.Stdout)

	// Exit with error code for non-2xx responses
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request returned %s", resp.Status)
	}

	return nil
}

// PrintEnhanced outputs only the response body with enhanced formatting.
func (f *Formatter) PrintEnhanced(resp *executor.Response) error {
	if len(resp.Body) > 0 {
		f.printEnhancedBody(resp.Body)
	}

	// Exit with error code for non-2xx responses
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request returned %s", resp.Status)
	}

	return nil
}

// PrintRaw outputs the response body as pretty-printed JSON.
func (f *Formatter) PrintRaw(resp *executor.Response) error {
	if len(resp.Body) > 0 {
		var jsonData interface{}
		if err := json.Unmarshal(resp.Body, &jsonData); err == nil {
			prettyBody, _ := json.MarshalIndent(jsonData, "", "  ")
			os.Stdout.Write(prettyBody)
			os.Stdout.WriteString("\n")
		} else {
			os.Stdout.Write(resp.Body)
		}
	}

	// Exit with error code for non-2xx responses
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request returned %s", resp.Status)
	}

	return nil
}

// PrintBodyJSON outputs only the body as JSON.
func (f *Formatter) PrintBodyJSON(resp *executor.Response) error {
	if len(resp.Body) > 0 {
		os.Stdout.Write(resp.Body)
		os.Stdout.WriteString("\n")
	}

	// Exit with error code for non-2xx responses
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request returned %s", resp.Status)
	}

	return nil
}

// PrintFullJSON outputs the full response as a single JSON object.
func (f *Formatter) PrintFullJSON(resp *executor.Response) error {
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

// printEnhancedBody displays the response body in an enhanced format.
// Objects are shown as properties, arrays as tables, with nested support.
func (f *Formatter) printEnhancedBody(body []byte) {
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		// Not JSON, print as-is
		fmt.Fprintf(os.Stdout, "%s\n", string(body))
		return
	}

	f.printValue(data, 0)
}

// printValue prints a value based on its type
func (f *Formatter) printValue(value interface{}, indent int) {
	switch v := value.(type) {
	case map[string]interface{}:
		f.printObject(v, indent)
	case []interface{}:
		f.printArray(v, indent)
	default:
		f.printPrimitive(v, indent)
	}
}

// printObject prints an object as properties with title-cased keys
func (f *Formatter) printObject(obj map[string]interface{}, indent int) {
	prefix := strings.Repeat("  ", indent)
	
	for key, value := range obj {
		titleKey := toTitleCase(key)
		fmt.Fprintf(os.Stdout, "%s%s: ", prefix, titleKey)
		
		switch v := value.(type) {
		case map[string]interface{}:
			fmt.Fprintln(os.Stdout)
			f.printObject(v, indent+1)
		case []interface{}:
			fmt.Fprintln(os.Stdout)
			f.printArray(v, indent+1)
		default:
			f.printPrimitiveInline(v)
			fmt.Fprintln(os.Stdout)
		}
	}
}

// printArray prints an array as a table or list
func (f *Formatter) printArray(arr []interface{}, indent int) {
	if len(arr) == 0 {
		fmt.Fprintln(os.Stdout, "(empty)")
		return
	}

	// Check if all elements are objects (table mode)
	allObjects := true
	for _, item := range arr {
		if _, ok := item.(map[string]interface{}); !ok {
			allObjects = false
			break
		}
	}

	if allObjects && len(arr) > 0 {
		f.printTable(arr, indent)
	} else {
		// List mode for primitives or mixed types
		prefix := strings.Repeat("  ", indent)
		for i, item := range arr {
			fmt.Fprintf(os.Stdout, "%s[%d] ", prefix, i)
			switch v := item.(type) {
			case map[string]interface{}:
				fmt.Fprintln(os.Stdout)
				f.printObject(v, indent+1)
			case []interface{}:
				fmt.Fprintln(os.Stdout)
				f.printArray(v, indent+1)
			default:
				f.printPrimitiveInline(v)
				fmt.Fprintln(os.Stdout)
			}
		}
	}
}

// printTable prints an array of objects as a formatted table
func (f *Formatter) printTable(arr []interface{}, indent int) {
	if len(arr) == 0 {
		return
	}

	// Collect all unique keys from all objects
	keySet := make(map[string]bool)
	for _, item := range arr {
		if obj, ok := item.(map[string]interface{}); ok {
			for key := range obj {
				keySet[key] = true
			}
		}
	}

	// Convert to sorted slice (for consistent column order)
	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Calculate column widths
	colWidths := make(map[string]int)
	for _, key := range keys {
		titleKey := toTitleCase(key)
		colWidths[key] = len(titleKey)
		
		// Check max width from values
		for _, item := range arr {
			if obj, ok := item.(map[string]interface{}); ok {
				if val, exists := obj[key]; exists {
					valStr := formatPrimitive(val)
					if len(valStr) > colWidths[key] {
						colWidths[key] = len(valStr)
					}
				}
			}
		}
	}

	// Print header
	prefix := strings.Repeat("  ", indent)
	header := prefix
	for _, key := range keys {
		titleKey := toTitleCase(key)
		header += fmt.Sprintf("%-"+fmt.Sprint(colWidths[key]+2)+"s", titleKey)
	}
	fmt.Fprintln(os.Stdout, header)

	// Print separator
	separator := prefix
	for _, key := range keys {
		separator += strings.Repeat("-", colWidths[key]+2)
	}
	fmt.Fprintln(os.Stdout, separator)

	// Print rows
	for _, item := range arr {
		if obj, ok := item.(map[string]interface{}); ok {
			row := prefix
			for _, key := range keys {
				val := ""
				if v, exists := obj[key]; exists {
					val = formatPrimitive(v)
				}
				row += fmt.Sprintf("%-"+fmt.Sprint(colWidths[key]+2)+"s", val)
			}
			fmt.Fprintln(os.Stdout, row)
		}
	}
}

// printPrimitive prints a primitive value with indentation
func (f *Formatter) printPrimitive(value interface{}, indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Fprintf(os.Stdout, "%s%s\n", prefix, formatPrimitive(value))
}

// printPrimitiveInline prints a primitive value without newline
func (f *Formatter) printPrimitiveInline(value interface{}) {
	fmt.Fprintf(os.Stdout, "%s", formatPrimitive(value))
}

// formatPrimitive formats a primitive value as string
func formatPrimitive(value interface{}) string {
	if value == nil {
		return "null"
	}
	
	switch v := value.(type) {
	case string:
		return v
	case float64:
		// Check if it's an integer
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// toTitleCase converts a string to Title Case
func toTitleCase(s string) string {
	if s == "" {
		return s
	}

	// Replace underscores and hyphens with spaces
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")

	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}

	return strings.Join(words, " ")
}
