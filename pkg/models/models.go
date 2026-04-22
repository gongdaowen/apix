package models

// APIOperation represents a single operation extracted from an OpenAPI spec.
type APIOperation struct {
	OperationID string
	Summary     string
	Description string
	Method      string
	Path        string
	Tags        []string
	Parameters  []ParameterSpec
	RequestBody *RequestBodySpec
	Security    []SecurityRequirement
}

// ParameterSpec describes a single parameter of an operation.
type ParameterSpec struct {
	Name        string
	In          string // "path", "query", "header"
	Required    bool
	Description string
	Schema      *SchemaInfo
}

// SchemaInfo holds simplified schema type information.
type SchemaInfo struct {
	Type    string
	Format  string
	Default interface{}
	Enum    []interface{}
}

// RequestBodySpec describes the request body of an operation.
type RequestBodySpec struct {
	Required    bool
	Description string
	ContentType string
	Schema      *SchemaInfo
	Example     interface{}
}

// SecurityRequirement represents a security scheme requirement.
type SecurityRequirement map[string][]string

// ResolvedRequest is a fully resolved request ready to be sent.
type ResolvedRequest struct {
	Operation *APIOperation
	BaseURL   string
	Method    string
	Path      string
	Query     map[string]string
	Headers   map[string]string
	Body      []byte
}

// SecurityScheme represents a security scheme from OpenAPI.
type SecurityScheme struct {
	Type   string // "apiKey", "http", "oauth2", "openIdConnect"
	Scheme string // "bearer", "basic"
	In     string // "header", "query", "cookie"
	Name   string // parameter name for apiKey
}
