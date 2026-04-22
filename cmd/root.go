package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gongdaowen/apix/internal/builder"
	"github.com/gongdaowen/apix/internal/executor"
	"github.com/gongdaowen/apix/internal/formatter"
	"github.com/gongdaowen/apix/internal/i18n"
	"github.com/gongdaowen/apix/internal/models"
	"github.com/gongdaowen/apix/internal/parser"
	"github.com/gongdaowen/apix/internal/resolver"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
)

var (
	specPath    string
	bodyFile    string
	headers     []string
	token       string
	apiKey      string
	envProfile  string
	serverIndex int
	baseUrl     string
	dryRun      bool
	outputFull  bool
	outputRaw   bool
	outputJSON  bool
	debugMode   bool
)

// Global translator instance
var translator *i18n.Translator

// Version is set at build time
var version = "dev"

var loadedSpec *openapi3.T
var globalSpecLoader *parser.SpecLoader

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "apix [command] [flags]",
	Short:   "", // Will be set in init() based on language
	Long:    "", // Will be set in init() based on language
	Example: "", // Will be set in init() based on language
	Version: "", // Will be set via SetVersion
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		// Update help text with server info if spec is loaded
		if loadedSpec != nil {
			servers := parser.GetServersInfo(loadedSpec)
			if len(servers) > 1 {
				cmd.Long = buildLongHelpWithServers(servers)
			} else {
				cmd.Long = buildLongHelp()
			}
		} else {
			cmd.Long = buildLongHelp()
		}
		cmd.Short = translator.T("root.short")
		cmd.Example = buildExamples()
		_ = cmd.Help()
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Set dynamic help text based on current language
		cmd.Short = translator.T("root.short")
		cmd.Long = buildLongHelp()
		cmd.Example = buildExamples()

		// Auto-detect spec file if not provided
		if specPath == "" && envProfile == "" {
			autoDetected := autoDetectSpecFile()
			if autoDetected != "" {
				specPath = autoDetected
			}
		}

		// Handle environment profile (-P flag)
		if envProfile != "" {
			specPath = findEnvProfile(envProfile)
			if specPath == "" {
				errMsg := translator.TF("error.env_profile_not_found", envProfile) + "\n\n" +
					translator.T("error.env_profile.hint")
				return fmt.Errorf("%s", errMsg)
			}
		}

		// If no spec provided at this point, we can't proceed
		if specPath == "" {
			errMsg := translator.T("error.no_spec") + "\n\n" +
				translator.T("error.no_spec.hint") + "\n\n" +
				"1. " + translator.T("error.no_spec.option1") + "\n" +
				"   apix --spec openapi.yaml <operation>\n" +
				"   apix -s https://api.example.com/openapi.json <operation>\n\n" +
				"2. " + translator.T("error.no_spec.option2") + "\n" +
				"   apix -P dev <operation>\n\n" +
				"3. " + translator.T("error.no_spec.option3") + "\n" +
				translator.T("error.no_spec.examples") + "\n" +
				"   apix --help\n" +
				"   apix getPet --petId 123\n" +
				"   apix -P prod listUsers\n\n" +
				translator.T("error.no_spec.more_info")
			return fmt.Errorf("%s", errMsg)
		}

		// Load the spec
		specLoader := parser.NewSpecLoader()
		doc, err := specLoader.Load(specPath)
		if err != nil {
			errMsg := translator.TF("error.load_spec_failed", specPath) + "\n\n" +
				translator.T("error.load_spec.check") + "\n" +
				"  - " + translator.T("error.load_spec.check.path") + "\n" +
				"  - " + translator.T("error.load_spec.check.exists") + "\n" +
				"  - " + translator.T("error.load_spec.check.valid")
			return fmt.Errorf("%s\n\n%w", errMsg, err)
		}
		loadedSpec = doc
		globalSpecLoader = specLoader

		return nil
	},
}

// createOperationRunner returns a RunE function for a specific operation
func createOperationRunner(doc *openapi3.T, op *models.APIOperation) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Collect parameter values from flags
		params := make(map[string]string)
		for _, param := range op.Parameters {
			value, _ := cmd.Flags().GetString(param.Name)
			if value != "" {
				params[param.Name] = value
			} else if param.Required {
				return fmt.Errorf("required parameter %q not provided", param.Name)
			}
		}

		// Convert params map to slice of "key=value" strings for resolver compatibility
		var paramSlice []string
		for k, v := range params {
			paramSlice = append(paramSlice, k+"="+v)
		}

		// Resolve parameters
		var baseURL string
		if baseUrl != "" {
			// Use explicitly provided base URL
			baseURL = baseUrl
		} else {
			// Auto-detect or use servers from spec
			baseURL = parser.GetBaseURL(doc, globalSpecLoader, serverIndex)
		}
		paramResolver := resolver.NewParamResolver(op, paramSlice, bodyFile, baseURL)
		resolved, err := paramResolver.Resolve()
		if err != nil {
			return err
		}

		// Build request
		reqBuilder := builder.NewRequestBuilder(resolved, headers, token, apiKey)
		httpReq, err := reqBuilder.Build()
		if err != nil {
			return fmt.Errorf("failed to build request: %w", err)
		}

		// Dry run mode
		if dryRun {
			fmtr := formatter.NewFormatter()
			// Get headers from request
			reqHeaders := make(map[string]string)
			for key, values := range httpReq.Header {
				if len(values) > 0 {
					reqHeaders[key] = values[0]
				}
			}
			// Get body from resolved request (before it was turned into ReadCloser)
			var bodyBytes []byte
			if resolved.Body != nil {
				bodyBytes, _ = json.Marshal(resolved.Body)
			}
			fmtr.PrintCurl(httpReq.Method, httpReq.URL.String(), reqHeaders, bodyBytes)
			return nil
		}

		// Execute
		exe := executor.NewExecutor(debugMode)
		resp, err := exe.Do(httpReq)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		// Format output based on mode
		fmtr := formatter.NewFormatter()
		
		// Handle --full --json combination: full response as JSON
		if outputFull && outputJSON {
			fmtr.SetMode(formatter.OutputFullJSON)
		} else if outputJSON {
			// --json only: body as JSON
			fmtr.SetMode(formatter.OutputBodyJSON)
		} else if outputFull {
			// --full only: complete response with enhanced body
			fmtr.SetMode(formatter.OutputFull)
		} else if outputRaw {
			// --raw only: body as pretty-printed JSON
			fmtr.SetMode(formatter.OutputRaw)
		}
		// Default: enhanced body only (OutputEnhanced)
		
		if err := fmtr.Print(resp); err != nil {
			return fmt.Errorf("failed to format response: %w", err)
		}

		return nil
	}
}

func init() {
	// Initialize translator based on environment variables
	translator = i18n.NewTranslator()

	// Set initial help text (will be updated in Execute based on detected language)
	rootCmd.Short = ""
	rootCmd.Long = ""
	rootCmd.Example = ""

	// Set version
	rootCmd.Version = version

	// Customize help template to reorder sections
	rootCmd.SetHelpTemplate(`{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)

	// Register flags with current language descriptions
	rootCmd.PersistentFlags().StringVarP(&specPath, "spec", "s", "", translator.T("flag.spec"))
	rootCmd.PersistentFlags().StringVarP(&envProfile, "profile", "P", "", translator.T("flag.profile"))
	rootCmd.PersistentFlags().StringVar(&baseUrl, "base-url", "", translator.T("flag.base_url"))
	rootCmd.PersistentFlags().IntVar(&serverIndex, "server", 0, translator.T("flag.server"))

	// Output format flags - must be PersistentFlags to work on subcommands
	rootCmd.PersistentFlags().BoolVar(&outputFull, "full", false, translator.T("flag.full"))
	rootCmd.PersistentFlags().BoolVar(&outputRaw, "raw", false, translator.T("flag.raw"))
	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, translator.T("flag.json"))
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, translator.T("flag.debug"))
}

// autoDetectSpecFile automatically detects OpenAPI spec files in the current directory
func autoDetectSpecFile() string {
	// Priority order for default filenames
	defaultFiles := []string{
		"openapi.yaml",
		"openapi.yml",
		"openapi.json",
		"api.yaml",
		"api.yml",
		"api.json",
		"swagger.yaml",
		"swagger.yml",
		"swagger.json",
	}

	for _, filename := range defaultFiles {
		if _, err := os.Stat(filename); err == nil {
			return filename
		}
	}

	return ""
}

// findEnvProfile finds spec file based on environment name
// e.g., "dev" -> "openapi-dev.yaml", "api-dev.yaml", etc.
func findEnvProfile(env string) string {
	// Try different naming patterns
	patterns := []string{
		fmt.Sprintf("openapi-%s.yaml", env),
		fmt.Sprintf("openapi-%s.yml", env),
		fmt.Sprintf("openapi-%s.json", env),
		fmt.Sprintf("api-%s.yaml", env),
		fmt.Sprintf("api-%s.yml", env),
		fmt.Sprintf("api-%s.json", env),
		fmt.Sprintf("swagger-%s.yaml", env),
		fmt.Sprintf("swagger-%s.yml", env),
		fmt.Sprintf("swagger-%s.json", env),
	}

	for _, pattern := range patterns {
		if _, err := os.Stat(pattern); err == nil {
			return pattern
		}
	}

	return ""
}

// buildLongHelp builds the long help text based on current language
func buildLongHelp() string {
	return fmt.Sprintf(`%s

%s:
  - %s
  - %s
  - %s
  - %s
  - %s
  - %s
  - %s`,
		translator.T("root.long"),
		"Features",
		translator.T("feature.load_specs"),
		translator.T("feature.auto_detect"),
		translator.T("feature.auto_generate"),
		translator.T("feature.support_params"),
		translator.T("feature.output_formats"),
		translator.T("feature.dry_run"),
		translator.T("feature.debug_mode"))
}

// buildLongHelpWithServers builds help text with server information
func buildLongHelpWithServers(servers []map[string]string) string {
	baseHelp := buildLongHelp()

	var serversSection strings.Builder
	serversSection.WriteString("\n\n")
	serversSection.WriteString(translator.T("op.servers") + "\n")
	for _, srv := range servers {
		desc := ""
		if srv["description"] != "" {
			desc = fmt.Sprintf(" - %s", srv["description"])
		}
		serversSection.WriteString(fmt.Sprintf("  [%s] %s%s\n", srv["index"], srv["url"], desc))
	}

	return baseHelp + serversSection.String()
}

// buildExamples builds examples based on current language
func buildExamples() string {
	return fmt.Sprintf(`  # %s
  apix

  # %s
  apix getPet --petId 123

  # %s
  apix -P dev getPet --petId 123

  # %s
  apix createPet -b pet.json -t YOUR_TOKEN

  # %s
  apix getPet --petId 123 --dry-run`,
		translator.T("example.load_spec"),
		translator.T("example.call_operation"),
		translator.T("example.use_env_profile"),
		translator.T("example.auth_request"),
		translator.T("example.preview_curl"))
}

// generateOperationHelp creates detailed help text for an operation
func generateOperationHelp(doc *openapi3.T, op models.APIOperation) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s\n\n", op.Description))
	sb.WriteString(fmt.Sprintf("%s %s\n", translator.T("op.operation_id"), op.OperationID))
	sb.WriteString(fmt.Sprintf("%s %s\n", translator.T("op.method"), op.Method))
	sb.WriteString(fmt.Sprintf("%s %s\n\n", translator.T("op.path"), op.Path))

	// Show available servers if multiple exist
	servers := parser.GetServersInfo(doc)
	if len(servers) > 1 {
		sb.WriteString(translator.T("op.servers") + "\n")
		for _, srv := range servers {
			defaultMarker := ""
			if srv["index"] == "0" {
				defaultMarker = " " + translator.T("op.default")
			}
			desc := ""
			if srv["description"] != "" {
				desc = fmt.Sprintf(" - %s", srv["description"])
			}
			sb.WriteString(fmt.Sprintf("  [%s] %s%s%s\n", srv["index"], srv["url"], desc, defaultMarker))
		}
		sb.WriteString(fmt.Sprintf("\n%s\n  --server <index>   %s\n\n", translator.T("op.server_flag_hint"), translator.T("flag.server")))
	}

	if len(op.Parameters) > 0 {
		sb.WriteString(translator.T("op.parameters") + "\n")
		for _, param := range op.Parameters {
			required := ""
			if param.Required {
				required = " " + translator.T("op.required")
			}
			sb.WriteString(fmt.Sprintf("  --%s: %s%s [%s]\n",
				param.Name,
				param.Description,
				required,
				param.In))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(translator.T("op.common_flags") + "\n")
	sb.WriteString("  -b, --body string    " + translator.T("flag.body") + "\n")
	sb.WriteString("  -H, --header strings " + translator.T("flag.header") + "\n")
	sb.WriteString("  -t, --token string   " + translator.T("flag.token") + "\n")
	sb.WriteString("  -k, --key string     " + translator.T("flag.key") + "\n")
	sb.WriteString("      --dry-run        " + translator.T("flag.dry_run") + "\n")
	sb.WriteString("      --raw            " + translator.T("flag.raw") + "\n")
	sb.WriteString("      --json           " + translator.T("flag.json") + "\n")
	sb.WriteString("      --debug          " + translator.T("flag.debug") + "\n")

	return sb.String()
}

// generateOperationExample creates example usage for an operation
func generateOperationExample(op models.APIOperation) string {
	var examples []string

	// Build a basic example with required parameters
	var requiredParams []string
	for _, param := range op.Parameters {
		if param.Required {
			exampleValue := getExampleValue(param)
			requiredParams = append(requiredParams, fmt.Sprintf("--%s %s", param.Name, exampleValue))
		}
	}

	basicExample := fmt.Sprintf("  # %s\n  apix %s %s",
		translator.TF("op.example.call", op.Summary),
		op.OperationID,
		strings.Join(requiredParams, " "))

	examples = append(examples, basicExample)

	// Add example with body if operation has request body
	if op.RequestBody != nil {
		bodyExample := fmt.Sprintf("  # %s\n  apix %s -b request.json %s",
			translator.TF("op.example.with_body", op.Summary),
			op.OperationID,
			strings.Join(requiredParams, " "))
		examples = append(examples, bodyExample)
	}

	// Add dry-run example
	dryRunExample := fmt.Sprintf("  # %s\n  apix %s %s --dry-run",
		translator.TF("op.example.preview", op.Summary),
		op.OperationID,
		strings.Join(requiredParams, " "))
	examples = append(examples, dryRunExample)

	return strings.Join(examples, "\n\n")
}

// getExampleValue returns an example value based on parameter type
func getExampleValue(param models.ParameterSpec) string {
	if param.Schema != nil {
		switch param.Schema.Type {
		case "integer", "number":
			return "1"
		case "boolean":
			return "true"
		}
	}
	// Fallback to parameter's own type field
	switch param.Type {
	case "integer", "number":
		return "1"
	case "boolean":
		return "true"
	default:
		return fmt.Sprintf("<%s>", param.Name)
	}
}

// Execute adds all commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	// Update help text with current language
	rootCmd.Short = translator.T("root.short")
	rootCmd.Long = buildLongHelp()
	rootCmd.Example = buildExamples()

	// First, parse the spec flag to dynamically register subcommands
	// We need to do this before Cobra's Execute() runs

	// Try to load spec early for dynamic subcommand registration
	// This is a workaround for Cobra's initialization cycle issue
	specFlag := ""
	args := os.Args[1:]

	// Check if --spec or -s is provided in command line args
	for i, arg := range args {
		if arg == "--spec" || arg == "-s" {
			if i+1 < len(args) {
				specFlag = args[i+1]
			}
		} else if strings.HasPrefix(arg, "--spec=") {
			specFlag = strings.TrimPrefix(arg, "--spec=")
		} else if strings.HasPrefix(arg, "-s=") {
			specFlag = strings.TrimPrefix(arg, "-s=")
		}
	}

	// Also check environment profile flag
	for i, arg := range args {
		if arg == "--profile" || arg == "-P" {
			if i+1 < len(args) {
				specFlag = findEnvProfile(args[i+1])
			}
		} else if strings.HasPrefix(arg, "--profile=") {
			specFlag = findEnvProfile(strings.TrimPrefix(arg, "--profile="))
		} else if strings.HasPrefix(arg, "-P=") {
			specFlag = findEnvProfile(strings.TrimPrefix(arg, "-P="))
		}
	}

	// If no spec provided via flags, try auto-detection
	if specFlag == "" {
		specFlag = autoDetectSpecFile()
	}

	if specFlag != "" {
		// Load spec and register dynamic subcommands
		specLoader := parser.NewSpecLoader()
		doc, err := specLoader.Load(specFlag)
		if err == nil {
			// Update root command help with server info if available
			servers := parser.GetServersInfo(doc)
			if len(servers) > 1 {
				rootCmd.Long = buildLongHelpWithServers(servers)
			}

			operations := parser.ExtractOperations(doc)
			for _, op := range operations {
				if op.OperationID == "" {
					continue // Skip operations without ID
				}

				// Create a new command for this operation
				opCmd := &cobra.Command{
					Use:     op.OperationID + " [flags]",
					Short:   op.Summary,
					Long:    generateOperationHelp(doc, op),
					Example: generateOperationExample(op),
					RunE:    createOperationRunner(doc, &op),
				}

				// Add flags for each parameter
				for _, param := range op.Parameters {
					flagName := param.Name
					if strings.ContainsAny(flagName, ".-") {
						flagName = strings.ReplaceAll(flagName, ".", "_")
						flagName = strings.ReplaceAll(flagName, "-", "_")
					}

					description := param.Description
					if param.Required {
						description += " (required)"
					}

					opCmd.Flags().StringP(param.Name, "", "", description)
				}

				// Add common flags with current language descriptions
				opCmd.Flags().StringVarP(&bodyFile, "body", "b", "", translator.T("flag.body"))
				opCmd.Flags().StringSliceVarP(&headers, "header", "H", nil, translator.T("flag.header"))
				opCmd.Flags().StringVarP(&token, "token", "t", "", translator.T("flag.token"))
				opCmd.Flags().StringVarP(&apiKey, "key", "k", "", translator.T("flag.key"))
				opCmd.Flags().BoolVar(&dryRun, "dry-run", false, translator.T("flag.dry_run"))

				rootCmd.AddCommand(opCmd)
			}
		}
	}

	return rootCmd.Execute()
}

// SetVersion sets the version string (called from main)
func SetVersion(v string) {
	version = v
	rootCmd.Version = version
}

// SetBuildInfo sets build time and commit hash (called from main)
func SetBuildInfo(buildTime, commitHash string) {
	// These can be used for extended version info if needed
}
