package i18n

func loadEnglishMessages() map[string]string {
	return map[string]string{
		// Root command
		"root.short": "Universal CLI tool that builds and sends HTTP requests from OpenAPI specs",
		"root.long": `apix reads an OpenAPI 3.x specification and lets you call any endpoint
from the command line. It supports interactive mode, direct operation selection,
and flexible parameter input.`,

		// Features
		"feature.load_specs":     "Load OpenAPI specs from local files or URLs",
		"feature.auto_detect":    "Auto-detect standard filenames in current directory (openapi.yaml, api.yaml, etc.)",
		"feature.auto_generate":  "Auto-generate CLI commands for each API operation",
		"feature.support_params": "Support for path parameters, query parameters, headers, and request bodies",
		"feature.output_formats": "Multiple output formats (pretty-printed, raw, JSON)",
		"feature.dry_run":        "Dry-run mode to preview curl commands",
		"feature.debug_mode":     "Debug mode to inspect request details",

		// Examples
		"example.load_spec":      "Auto-detect spec and list available operations",
		"example.call_operation": "Call a specific API operation",
		"example.use_env_profile": "Use environment profile (dev/prod/staging)",
		"example.auth_request":   "Send request with authentication",
		"example.preview_curl":   "Preview curl command without sending",

		// Flags - Global
		"flag.spec":    "Path or URL to OpenAPI specification file",
		"flag.profile": "Use environment name to auto-find spec file (e.g., dev -> openapi-dev.yaml)",
		"flag.base_url": "Override server URL (takes precedence over servers defined in the spec)",
		"flag.server":  "Select server index (when multiple servers are defined), default is 0 (first server)",
		"flag.lang":    "Set language preference (en/zh). Auto-detects if not specified",
		"flag.raw":     "Output only response body without headers or status",
		"flag.json":    "Output full response as JSON including status, headers, and body",
		"flag.debug":   "Show debug information including request URL and method",

		// Flags - Operation
		"flag.body":    "JSON file for request body",
		"flag.header":  "Header as 'Key: Value' (can be specified multiple times)",
		"flag.token":   "Bearer token for authentication",
		"flag.key":     "API key for authentication",
		"flag.dry_run": "Print curl command without sending request",

		// Error messages
		"error.no_spec":           "OpenAPI specification is required but not provided",
		"error.no_spec.hint":      "You must provide an OpenAPI spec in one of these ways:",
		"error.no_spec.option1":   "Use --spec flag with a file path or URL:",
		"error.no_spec.option2":   "Use environment profile (-P dev/prod/staging):",
		"error.no_spec.option3":   "Place spec file in current directory with standard naming:",
		"error.no_spec.examples":  "Examples:",
		"error.no_spec.more_info": "Run 'apix --help' for more information.",

		"error.load_spec_failed":       "Failed to load OpenAPI specification from %q",
		"error.load_spec.check":        "Please check:",
		"error.load_spec.check.path":   "The file path or URL is correct",
		"error.load_spec.check.exists": "The file exists and is accessible",
		"error.load_spec.check.valid":  "The file contains a valid OpenAPI 3.x specification",

		"error.env_profile_not_found": "No spec file found for environment %q",
		"error.env_profile.hint":      "Hint: Ensure openapi-{env}.yaml, api-{env}.yaml, or swagger-{env}.yaml exists",

		// Help text
		"help.usage":              "Usage:",
		"help.examples":           "Examples:",
		"help.flags":              "Flags:",
		"help.global_flags":       "Global Flags:",
		"help.available_commands": "Available Commands:",
		"help.more_info":          "Use \"%s [command] --help\" for more information about a command.",

		// Operation help
		"op.operation_id":      "Operation ID:",
		"op.method":            "Method:",
		"op.path":              "Path:",
		"op.servers":           "Available Servers:",
		"op.default":           "(default)",
		"op.server_flag_hint":  "Use --server flag to select a server:",
		"op.parameters":        "Parameters:",
		"op.common_flags":      "Common Flags:",
		"op.required":          "(required)",
		"op.example.call":      "Call %s",
		"op.example.with_body": "Call %s with request body",
		"op.example.preview":   "Preview %s request",
	}
}
