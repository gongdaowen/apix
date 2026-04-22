package parser

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/apix-cli/apix/internal/models"
	"github.com/getkin/kin-openapi/openapi3"
)

// SpecLoader loads and parses OpenAPI 3.x specifications.
type SpecLoader struct {
	specPath string // Store the spec path for auto-detection
}

func NewSpecLoader() *SpecLoader {
	return &SpecLoader{}
}

// SetSpecPath sets the spec path for auto-detection
func (l *SpecLoader) SetSpecPath(path string) {
	l.specPath = path
}

// Load reads an OpenAPI spec from a file path or URL.
func (l *SpecLoader) Load(pathOrURL string) (*openapi3.T, error) {
	// Store the spec path for auto-detection
	l.specPath = pathOrURL
	
	loader := openapi3.NewLoader()

	// Check if it's a URL
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		url, err := url.Parse(pathOrURL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}
		return loader.LoadFromURI(url)
	}

	// Read from local file
	data, err := os.ReadFile(pathOrURL)
	if err != nil {
		return nil, fmt.Errorf("cannot read spec file: %w", err)
	}

	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("cannot parse OpenAPI spec: %w", err)
	}

	// Validate with warning instead of error for compatibility issues
	if err := doc.Validate(context.Background()); err != nil {
		// Log warning but continue loading for CLI usage
		fmt.Fprintf(os.Stderr, "Warning: OpenAPI validation issue (continuing anyway): %v\n", err)
	}

	return doc, nil
}

// ExtractOperations extracts all operations from the OpenAPI document.
func ExtractOperations(doc *openapi3.T) []models.APIOperation {
	var operations []models.APIOperation

	for path, pathItem := range doc.Paths.Map() {
		// Extract path parameters
		pathParams := extractParameters(pathItem.Parameters, "path")

		// Process each HTTP method
		methods := map[string]*openapi3.Operation{
			"GET":    pathItem.Get,
			"POST":   pathItem.Post,
			"PUT":    pathItem.Put,
			"PATCH":  pathItem.Patch,
			"DELETE": pathItem.Delete,
			"HEAD":   pathItem.Head,
			"OPTIONS": pathItem.Options,
		}

		for method, op := range methods {
			if op == nil {
				continue
			}

			// Merge path-level parameters with operation-level parameters
			allParams := extractParameters(op.Parameters, "")
			allParams = append(allParams, pathParams...)

			operation := models.APIOperation{
				OperationID: op.OperationID,
				Summary:     op.Summary,
				Description: op.Description,
				Method:      method,
				Path:        path,
				Tags:        op.Tags,
				Parameters:  allParams,
				Security:    extractSecurityOrNil(op.Security),
			}

			// Extract request body
			if op.RequestBody != nil {
				operation.RequestBody = extractRequestBody(op.RequestBody.Value)
			}

			operations = append(operations, operation)
		}
	}

	return operations
}

// extractParameters converts OpenAPI parameter references to our models.
func extractParameters(params openapi3.Parameters, inFilter string) []models.ParameterSpec {
	var result []models.ParameterSpec

	for _, paramRef := range params {
		param := paramRef.Value
		if param == nil {
			continue
		}

		// Apply filter if specified (e.g., only "path" params)
		if inFilter != "" && param.In != inFilter {
			continue
		}

		spec := models.ParameterSpec{
			Name:        param.Name,
			In:          param.In,
			Required:    param.Required,
			Description: param.Description,
			Type:        "string", // Default to string for CLI
		}

		if param.Schema != nil && param.Schema.Value != nil {
			schemaInfo := extractSchema(param.Schema.Value)
			if schemaInfo != nil {
				spec.Type = schemaInfo.Type
			}
			spec.Schema = schemaInfo
		}

		result = append(result, spec)
	}

	return result
}

// extractRequestBody extracts request body information.
func extractRequestBody(body *openapi3.RequestBody) *models.RequestBodySpec {
	if body == nil || len(body.Content) == 0 {
		return nil
	}

	// Prefer JSON, fallback to first available
	var mediaType string
	var media *openapi3.MediaType
	for mt, m := range body.Content {
		if strings.Contains(mt, "json") {
			mediaType = mt
			media = m
			break
		}
	}
	if media == nil {
		for mt, m := range body.Content {
			mediaType = mt
			media = m
			break
		}
	}

	spec := &models.RequestBodySpec{
		Required:    body.Required,
		Description: body.Description,
		Content: map[string]interface{}{
			mediaType: map[string]interface{}{
				"schema":   extractSchema(media.Schema.Value),
				"example":  media.Example,
				"examples": media.Examples,
			},
		},
	}

	return spec
}

// extractSchema simplifies OpenAPI schema to our model.
func extractSchema(schema *openapi3.Schema) *models.SchemaInfo {
	if schema == nil {
		return nil
	}

	info := &models.SchemaInfo{
		Type:    "string", // Default to string for CLI simplicity
		Format:  schema.Format,
		Default: schema.Default,
	}

	if schema.Type != nil && len(*schema.Type) > 0 {
		info.Type = (*schema.Type)[0]
	}

	for _, v := range schema.Enum {
		info.Enum = append(info.Enum, v)
	}

	return info
}

// extractSecurityOrNil safely extracts security requirements, handling nil.
func extractSecurityOrNil(security *openapi3.SecurityRequirements) []interface{} {
	if security == nil {
		return nil
	}
	return extractSecurity(*security)
}

// extractSecurity converts OpenAPI security requirements.
func extractSecurity(security openapi3.SecurityRequirements) []interface{} {
	var result []interface{}
	for _, req := range security {
		result = append(result, req)
	}
	return result
}

// GetSecuritySchemes extracts all security schemes from the document.
func GetSecuritySchemes(doc *openapi3.T) map[string]interface{} {
	result := make(map[string]interface{})

	for name, schemeRef := range doc.Components.SecuritySchemes {
		if schemeRef == nil || schemeRef.Value == nil {
			continue
		}
		scheme := schemeRef.Value
		result[name] = map[string]interface{}{
			"type":   scheme.Type,
			"scheme": scheme.Scheme,
			"in":     scheme.In,
			"name":   scheme.Name,
		}
	}

	return result
}

// GetBaseURL extracts the server URL from the OpenAPI document.
// If index is provided and valid, uses that server; otherwise uses the first one.
// If no servers are defined, auto-detects based on spec path or defaults to localhost.
func GetBaseURL(doc *openapi3.T, loader *SpecLoader, index ...int) string {
	// If servers are defined, use them
	if len(doc.Servers) > 0 {
		// Use specified index if provided and valid
		if len(index) > 0 && index[0] >= 0 && index[0] < len(doc.Servers) {
			return doc.Servers[index[0]].URL
		}
		// Default to first server
		return doc.Servers[0].URL
	}

	// No servers defined, auto-detect
	if loader != nil && loader.specPath != "" {
		// Check if spec path is a URL
		if strings.HasPrefix(loader.specPath, "http://") || strings.HasPrefix(loader.specPath, "https://") {
			// Extract base URL from spec URL
			parsedURL, err := url.Parse(loader.specPath)
			if err == nil {
				// Construct base URL: scheme://host[:port]
				baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
				fmt.Fprintf(os.Stderr, "Info: No servers defined in spec, auto-detected from spec URL: %s\n", baseURL)
				return baseURL
			}
		}
	}

	// Default to localhost
	defaultURL := "http://localhost"
	fmt.Fprintf(os.Stderr, "Info: No servers defined in spec, using default: %s\n", defaultURL)
	return defaultURL
}

// GetServersInfo returns information about all available servers
func GetServersInfo(doc *openapi3.T) []map[string]string {
	var servers []map[string]string
	for i, server := range doc.Servers {
		info := map[string]string{
			"index":       fmt.Sprintf("%d", i),
			"url":         server.URL,
			"description": server.Description,
		}
		servers = append(servers, info)
	}
	return servers
}
