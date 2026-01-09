package commands

import (
	"os"

	"github.com/miguel-bandeira-infosistema/schemactl/internal/client"
	"github.com/miguel-bandeira-infosistema/schemactl/internal/config"
	"github.com/miguel-bandeira-infosistema/schemactl/internal/output"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version = "dev"

	// Global flags
	apiURL  string
	apiKey  string
	jsonOut bool

	// Shared client
	apiClient *client.Client
	cfg       *config.Config
)

var rootCmd = &cobra.Command{
	Use:     "schemactl",
	Aliases: []string{"sr", "sreg"},
	Short:   "CLI for Schema Registry API",
	Long: `schemactl is a command-line interface for interacting with the Schema Registry API.

It allows you to manage schemas, document types, countries, and query reference data.

Default API URL: https://api.docdigitizer.com/registry

Environment variables:
  SCHEMACTL_API_URL    API base URL (overrides default)
  SCHEMACTL_API_KEY    API key for authentication (optional)
  SCHEMACTL_TIMEOUT    Request timeout in seconds (default: 30)

Run without arguments to enter interactive shell mode.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand given, start interactive shell
		return RunShell()
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization for help and version commands
		if cmd.Name() == "help" || cmd.Name() == "version" {
			return nil
		}

		// Load configuration
		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}

		// Apply command-line overrides
		cfg.WithAPIURL(apiURL).WithAPIKey(apiKey)

		// Validate configuration
		if err := cfg.Validate(); err != nil {
			return err
		}

		// Set JSON output mode
		output.JSONOutput = jsonOut

		// Create API client
		apiClient = client.New(cfg)

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "API base URL (env: SCHEMACTL_API_URL)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key (env: SCHEMACTL_API_KEY)")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Output as JSON")

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(schemasCmd)
	rootCmd.AddCommand(docTypesCmd)
	rootCmd.AddCommand(countriesCmd)
	rootCmd.AddCommand(referenceDataCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintSuccess("schemactl version " + Version)
	},
}

// GetClient returns the API client (for use by subcommands)
func GetClient() *client.Client {
	return apiClient
}

// ExitOnError prints an error and exits with code 1
func ExitOnError(err error) {
	if err != nil {
		output.PrintError(err)
		os.Exit(1)
	}
}
