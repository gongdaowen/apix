package builder

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gongdaowen/apix/internal/models"
)

// RequestBuilder constructs an HTTP request from a resolved request.
type RequestBuilder struct {
	resolved *models.ResolvedRequest
	headers  []string
	token    string
	apiKey   string
}

func NewRequestBuilder(resolved *models.ResolvedRequest, headers []string, token, apiKey string) *RequestBuilder {
	return &RequestBuilder{
		resolved: resolved,
		headers:  headers,
		token:    token,
		apiKey:   apiKey,
	}
}

// Build constructs and returns the final *http.Request.
func (b *RequestBuilder) Build() (*http.Request, error) {
	// Build full URL
	fullURL := b.buildURL()

	// Create request
	var req *http.Request
	var err error

	if b.resolved.Body != nil {
		var bodyBytes []byte
		switch v := b.resolved.Body.(type) {
		case []byte:
			bodyBytes = v
		case string:
			bodyBytes = []byte(v)
		default:
			bodyBytes, err = json.Marshal(v)
			if err != nil {
				return nil, err
			}
		}
		req, err = http.NewRequest(b.resolved.Method, fullURL, bytes.NewReader(bodyBytes))
	} else {
		req, err = http.NewRequest(b.resolved.Method, fullURL, nil)
	}
	if err != nil {
		return nil, err
	}

	// Set headers from resolved params
	for key, value := range b.resolved.Headers {
		req.Header.Set(key, value)
	}

	// Set custom headers from --header flag
	for _, h := range b.headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	// Set Content-Type if body is present
	if b.resolved.Body != nil {
		if req.Header.Get("Content-Type") == "" {
			contentType := "application/json"
			if b.resolved.Operation.RequestBody != nil && len(b.resolved.Operation.RequestBody.Content) > 0 {
				// Get first content type from the Content map
				for ct := range b.resolved.Operation.RequestBody.Content {
					contentType = ct
					break
				}
			}
			req.Header.Set("Content-Type", contentType)
		}
	}

	// Handle authentication
	b.applyAuth(req)

	return req, nil
}

func (b *RequestBuilder) buildURL() string {
	baseURL := b.resolved.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost"
	}

	// Ensure no double slashes
	baseURL = strings.TrimSuffix(baseURL, "/")

	rawPath := b.resolved.Path

	// Append query parameters
	if len(b.resolved.Query) > 0 {
		params := url.Values{}
		for key, value := range b.resolved.Query {
			params.Set(key, value)
		}
		rawPath = rawPath + "?" + params.Encode()
	}

	return baseURL + rawPath
}

func (b *RequestBuilder) applyAuth(req *http.Request) {
	// Explicit --token flag
	if b.token != "" {
		req.Header.Set("Authorization", "Bearer "+b.token)
		return
	}

	// Explicit --api-key flag
	if b.apiKey != "" {
		req.Header.Set("X-API-Key", b.apiKey)
		return
	}

	// Auto-detect from OpenAPI security requirements
	for _, sec := range b.resolved.Operation.Security {
		if secMap, ok := sec.(map[string][]string); ok {
			for schemeName := range secMap {
				// These would come from parsed security schemes
				// For now, check env vars or prompt
				if envToken := getEnvAuth(schemeName); envToken != "" {
					req.Header.Set("Authorization", "Bearer "+envToken)
					return
				}
			}
		}
	}
}

func getEnvAuth(schemeName string) string {
	// Check environment variables like APIX_TOKEN_<SCHEME>
	// Simplified for now
	return ""
}
