package models

// ParameterSpec represents an OpenAPI parameter.
type ParameterSpec struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Required    bool        `json:"required"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Schema      *SchemaInfo `json:"schema,omitempty"`
}

// SchemaInfo represents simplified schema information.
type SchemaInfo struct {
	Type        string                 `json:"type"`
	Format      string                 `json:"format,omitempty"`
	Default     interface{}            `json:"default,omitempty"`
	Enum        []interface{}          `json:"enum,omitempty"`
	Properties  map[string]*SchemaInfo `json:"properties,omitempty"`
	Items       *SchemaInfo            `json:"items,omitempty"`
	Required    []string               `json:"required,omitempty"`
}

// RequestBodySpec represents an OpenAPI request body.
type RequestBodySpec struct {
	Required    bool                   `json:"required"`
	Description string                 `json:"description"`
	Content     map[string]interface{} `json:"content"`
}

// APIOperation represents a single operation extracted from an OpenAPI spec.
type APIOperation struct {
	OperationID   string           `json:"operationId"`
	Summary       string           `json:"summary"`
	Description   string           `json:"description"`
	Method        string           `json:"method"`
	Path          string           `json:"path"`
	Tags          []string         `json:"tags"`
	Parameters    []ParameterSpec  `json:"parameters"`
	RequestBody   *RequestBodySpec `json:"requestBody,omitempty"`
	Security      []interface{}    `json:"security"`
}