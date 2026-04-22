package resolver

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/apix-cli/apix/internal/models"
	"github.com/apix-cli/apix/internal/parser"
	"github.com/getkin/kin-openapi/openapi3"
)

// OperationResolver handles selecting an operation from the spec.
type OperationResolver struct {
	doc        *openapi3.T
	operations []models.APIOperation
}

func NewOperationResolver(doc *openapi3.T) *OperationResolver {
	return &OperationResolver{
		doc:        doc,
		operations: parser.ExtractOperations(doc),
	}
}

// Resolve returns the target operation. If operationID is empty, starts interactive selection.
func (r *OperationResolver) Resolve(operationID string) (*models.APIOperation, error) {
	if operationID != "" {
		return r.findByID(operationID)
	}
	return r.interactiveSelect()
}

func (r *OperationResolver) findByID(id string) (*models.APIOperation, error) {
	for _, op := range r.operations {
		if op.OperationID == id {
			return &op, nil
		}
	}
	return nil, fmt.Errorf("operation '%s' not found. Available operations:\n%s",
		id, r.formatOperationList())
}

func (r *OperationResolver) interactiveSelect() (*models.APIOperation, error) {
	if len(r.operations) == 0 {
		return nil, fmt.Errorf("no operations found in the OpenAPI spec")
	}

	// Build display options
	options := make([]string, len(r.operations))
	for i, op := range r.operations {
		options[i] = formatOperation(&op)
	}

	var selection string
	prompt := &survey.Select{
		Message: "Choose an API operation:",
		Options: options,
		PageSize: 15,
		Filter: func(filter string, value string, index int) bool {
			return strings.Contains(strings.ToLower(value), strings.ToLower(filter))
		},
	}

	if err := survey.AskOne(prompt, &selection); err != nil {
		return nil, fmt.Errorf("operation selection cancelled")
	}

	// Find selected operation
	for _, op := range r.operations {
		if formatOperation(&op) == selection {
			return &op, nil
		}
	}

	return nil, fmt.Errorf("failed to find selected operation")
}

func (r *OperationResolver) formatOperationList() string {
	var lines []string
	for _, op := range r.operations {
		lines = append(lines, fmt.Sprintf("  %s %s (%s)",
			op.Method, op.Path, op.OperationID))
	}
	return strings.Join(lines, "\n")
}

func formatOperation(op *models.APIOperation) string {
	label := fmt.Sprintf("%s %s", op.Method, op.Path)
	if op.OperationID != "" {
		label += " — " + op.OperationID
	}
	if op.Summary != "" {
		label += " | " + op.Summary
	}
	return label
}

// ParamResolver handles resolving parameters from CLI, file, or interactive input.
type ParamResolver struct {
	operation *models.APIOperation
	params    []string
	bodyFile  string
	baseURL   string
}

func NewParamResolver(operation *models.APIOperation, params []string, bodyFile string, baseURL string) *ParamResolver {
	return &ParamResolver{
		operation: operation,
		params:    params,
		bodyFile:  bodyFile,
		baseURL:   baseURL,
	}
}

// Resolve builds a ResolvedRequest by collecting all parameters.
func (r *ParamResolver) Resolve() (*models.ResolvedRequest, error) {
	cliParams := parseKeyValuePairs(r.params)

	resolved := &models.ResolvedRequest{
		Operation: r.operation,
		BaseURL:   r.baseURL,
		Method:    r.operation.Method,
		Path:      r.operation.Path,
		Query:     make(map[string]string),
		Headers:   make(map[string]string),
	}

	// Resolve parameters
	for _, param := range r.operation.Parameters {
		// Check CLI params first
		if val, ok := cliParams[param.Name]; ok {
			r.setParam(resolved, param, val)
			continue
		}

		// Required but not provided → interactive
		if param.Required {
			val, err := r.promptParam(&param)
			if err != nil {
				return nil, fmt.Errorf("parameter '%s': %w", param.Name, err)
			}
			r.setParam(resolved, param, val)
		}
	}

	// Resolve path parameters in URL template
	resolved.Path = replacePathParams(resolved.Path, resolved.Path, cliParams)

	// Resolve request body
	if r.operation.RequestBody != nil {
		body, err := r.resolveBody()
		if err != nil {
			return nil, fmt.Errorf("request body: %w", err)
		}
		resolved.Body = body
	}

	return resolved, nil
}

func (r *ParamResolver) setParam(resolved *models.ResolvedRequest, param models.ParameterSpec, value string) {
	switch param.In {
	case "path":
		resolved.Path = strings.ReplaceAll(resolved.Path, "{"+param.Name+"}", value)
	case "query":
		resolved.Query[param.Name] = value
	case "header":
		resolved.Headers[param.Name] = value
	}
}

func (r *ParamResolver) promptParam(param *models.ParameterSpec) (string, error) {
	var answer string
	prompt := &survey.Input{
		Message: fmt.Sprintf("%s (%s, required)", param.Name, param.In),
		Default: fmt.Sprintf("%v", getDefaultValue(param)),
	}

	if param.Description != "" {
		prompt.Help = param.Description
	}

	if err := survey.AskOne(prompt, &answer); err != nil {
		return "", err
	}
	return answer, nil
}

func (r *ParamResolver) resolveBody() (interface{}, error) {
	// 1. Try --body-file
	if r.bodyFile != "" {
		data, err := os.ReadFile(r.bodyFile)
		if err != nil {
			return nil, err
		}
		// Try to parse as JSON, otherwise return as string
		var jsonBody interface{}
		if err := json.Unmarshal(data, &jsonBody); err == nil {
			return jsonBody, nil
		}
		return string(data), nil
	}

	// 2. Try body from CLI params (--param body='{...}')
	// (Handled by checking if there's a special "body" param)

	// 3. Interactive: show schema hint and prompt for JSON
	if r.operation.RequestBody.Required {
		var answer string
		prompt := &survey.Input{
			Message: "Request body (JSON)",
			Help:    r.getBodyHint(),
		}
		if err := survey.AskOne(prompt, &answer); err != nil {
			return nil, err
		}
		// Try to parse as JSON
		var jsonBody interface{}
		if err := json.Unmarshal([]byte(answer), &jsonBody); err == nil {
			return jsonBody, nil
		}
		return answer, nil
	}

	return nil, nil
}

func (r *ParamResolver) getBodyHint() string {
	if r.operation.RequestBody == nil {
		return ""
	}
	rb := r.operation.RequestBody
	if rb != nil && len(rb.Content) > 0 {
		for _, content := range rb.Content {
			if contentMap, ok := content.(map[string]interface{}); ok {
				if example := contentMap["example"]; example != nil {
					return fmt.Sprintf("Example: %v", example)
				}
				if schema := contentMap["schema"]; schema != nil {
					if schemaInfo, ok := schema.(*models.SchemaInfo); ok {
						return fmt.Sprintf("Type: %s", schemaInfo.Type)
					}
				}
			}
		}
	}
	return "JSON object"
}

// parseKeyValuePairs parses ["key=value", "foo=bar"] into a map.
func parseKeyValuePairs(pairs []string) map[string]string {
	result := make(map[string]string)
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

// replacePathParams replaces {param} placeholders in path template.
func replacePathParams(template string, current string, params map[string]string) string {
	result := current
	for key, value := range params {
		result = strings.ReplaceAll(result, "{"+key+"}", value)
	}
	return result
}

func getDefaultValue(param *models.ParameterSpec) interface{} {
	if param.Schema != nil && param.Schema.Default != nil {
		return param.Schema.Default
	}
	switch param.Schema.Type {
	case "integer", "number":
		return 0
	case "boolean":
		return false
	case "array":
		return "[]"
	default:
		return ""
	}
}

// PromptForAuth interactively asks for missing auth credentials.
func PromptForAuth(schemes map[string]interface{}) (string, string) {
	var token, apiKey string

	for name, scheme := range schemes {
		if schemeMap, ok := scheme.(map[string]interface{}); ok {
			if schemeType, ok := schemeMap["type"].(string); ok {
				switch schemeType {
				case "http":
					if schemeScheme, ok := schemeMap["scheme"].(string); ok && schemeScheme == "bearer" {
						if token == "" {
							survey.AskOne(&survey.Password{
								Message: fmt.Sprintf("Bearer token for '%s':", name),
							}, &token)
						}
					}
				case "apiKey":
					if schemeIn, ok := schemeMap["in"].(string); ok {
						if apiKey == "" {
							survey.AskOne(&survey.Password{
								Message: fmt.Sprintf("API key for '%s' (in %s):", name, schemeIn),
							}, &apiKey)
						}
					}
				}
			}
		}
	}

	return token, apiKey
}

// FormatValue converts a value to string based on schema type.
func FormatValue(value interface{}, schemaType string) string {
	if value == nil {
		return ""
	}
	switch schemaType {
	case "integer", "number":
		return fmt.Sprintf("%v", value)
	case "boolean":
		b, _ := strconv.ParseBool(fmt.Sprintf("%v", value))
		return strconv.FormatBool(b)
	default:
		return fmt.Sprintf("%v", value)
	}
}
