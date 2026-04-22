package models

// ResolvedRequest holds all the information needed to build an HTTP request.
type ResolvedRequest struct {
	Operation   *APIOperation
	BaseURL     string
	Method      string
	Path        string
	Query       map[string]string
	Headers     map[string]string
	PathParams  map[string]string
	Body        interface{}
	Security    []interface{}
}