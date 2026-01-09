package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/miguel-bandeira-infosistema/schemactl/internal/client"
	"github.com/miguel-bandeira-infosistema/schemactl/internal/config"
	"github.com/miguel-bandeira-infosistema/schemactl/internal/output"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start interactive shell",
	Long:  "Start an interactive shell session where you can run commands without the 'schemactl' prefix",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunShell()
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}

// RunShell starts the interactive shell
func RunShell() error {
	reader := bufio.NewReader(os.Stdin)

	printBanner()
	fmt.Println()

	for {
		fmt.Print("schemactl> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nGoodbye!")
				return nil
			}
			return err
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Handle exit commands
		if input == "exit" || input == "quit" || input == "q" {
			fmt.Println("Goodbye!")
			return nil
		}

		// Handle help command specially
		if input == "help" || input == "?" {
			printShellHelp()
			continue
		}

		// Parse input into args
		args := parseArgs(input)
		if len(args) == 0 {
			continue
		}

		// Execute the command
		executeShellCommand(args)
	}
}

// parseArgs splits input string into arguments, respecting quotes
func parseArgs(input string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range input {
		switch {
		case r == '"' || r == '\'':
			if inQuote && r == quoteChar {
				inQuote = false
				quoteChar = 0
			} else if !inQuote {
				inQuote = true
				quoteChar = r
			} else {
				current.WriteRune(r)
			}
		case r == ' ' && !inQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

// executeShellCommand runs a command with the given args
func executeShellCommand(args []string) {
	// Create a fresh command tree for each execution
	cmd := buildRootCommand()
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println()
}

// buildRootCommand creates a new root command with all subcommands (except shell)
func buildRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "schemactl",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			if c.Name() == "help" || c.Name() == "version" {
				return nil
			}
			return initClient()
		},
	}

	// Add global flags
	cmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "API base URL")
	cmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key")
	cmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Output as JSON")

	// Add all commands
	cmd.AddCommand(buildVersionCmd())
	cmd.AddCommand(buildHealthCmd())
	cmd.AddCommand(buildSchemasCmd())
	cmd.AddCommand(buildDocTypesCmd())
	cmd.AddCommand(buildCountriesCmd())
	cmd.AddCommand(buildReferenceDataCmd())

	return cmd
}

func buildVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("schemactl version " + Version)
		},
	}
}

func buildHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check API health",
		Long:  "Check the health status of the API server and database connection",
		RunE: func(c *cobra.Command, args []string) error {
			health, err := GetClient().Health()
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(health)
			}
			statusIcon := "OK"
			if health.Status != "ok" {
				statusIcon = "ERROR"
			}
			dbIcon := "Connected"
			if health.Database != "connected" {
				dbIcon = "Disconnected"
			}
			fmt.Printf("Status:    %s\n", statusIcon)
			fmt.Printf("Database:  %s\n", dbIcon)
			fmt.Printf("Timestamp: %s\n", health.Timestamp.Format("2006-01-02 15:04:05"))
			return nil
		},
	}
}

func buildSchemasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schemas",
		Short: "Manage schemas",
		Long:  "Commands for managing JSON schemas in the registry",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List schemas",
		Long:  "List schemas with optional filtering",
		RunE: func(c *cobra.Command, args []string) error {
			opts := &client.ListSchemasOptions{}
			if s, _ := c.Flags().GetString("status"); s != "" {
				status := client.Status(s)
				opts.Status = &status
			}
			if s, _ := c.Flags().GetString("doc-type"); s != "" {
				opts.DocType = &s
			}
			if s, _ := c.Flags().GetString("country"); s != "" {
				opts.Country = &s
			}
			if s, _ := c.Flags().GetString("visibility"); s != "" {
				vis := client.Visibility(s)
				opts.Visibility = &vis
			}
			if s, _ := c.Flags().GetString("customer-id"); s != "" {
				opts.CustomerID = &s
			}
			opts.Limit, _ = c.Flags().GetInt("limit")
			opts.Offset, _ = c.Flags().GetInt("offset")
			result, err := GetClient().ListSchemas(opts)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(result)
			}
			headers := []string{"ID", "VERSION ID", "NAME", "DOC TYPE", "COUNTRY", "STATUS", "VER", "VISIBILITY"}
			rows := make([][]string, len(result.Data))
			for i, s := range result.Data {
				rows[i] = []string{
					s.PublicID,
					s.PublicVersionID,
					output.Truncate(s.Name, 30),
					s.DocTypeCode,
					output.PtrString(s.CountryCode, "-"),
					string(s.Status),
					strconv.Itoa(s.Version),
					string(s.Visibility),
				}
			}
			output.PrintTable(headers, rows)
			if result.Pagination.HasMore {
				fmt.Printf("\nShowing %d of %d schemas (use --offset to see more)\n",
					len(result.Data), result.Pagination.Total)
			}
			return nil
		},
	}
	listCmd.Flags().String("status", "", "Filter by status (draft, active, deprecated)")
	listCmd.Flags().StringP("doc-type", "t", "", "Filter by doc type code")
	listCmd.Flags().StringP("country", "c", "", "Filter by country code")
	listCmd.Flags().StringP("visibility", "v", "", "Filter by visibility (public, community, private)")
	listCmd.Flags().String("customer-id", "", "Filter by customer ID")
	listCmd.Flags().Int("limit", 50, "Number of items to return")
	listCmd.Flags().Int("offset", 0, "Number of items to skip")

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a schema",
		Long:  "Get a schema by publicId (sch_xxx) or publicVersionId (schv_xxx)",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			schema, err := GetClient().GetSchema(id)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(schema)
			}
			printSchemaDetails(schema)
			return nil
		},
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a schema",
		Long:  "Create a new schema in draft status",
		RunE: func(c *cobra.Command, args []string) error {
			name, _ := c.Flags().GetString("name")
			docType, _ := c.Flags().GetString("doc-type")
			contentStr, _ := c.Flags().GetString("content")
			description, _ := c.Flags().GetString("description")
			country, _ := c.Flags().GetString("country")
			visibility, _ := c.Flags().GetString("visibility")
			schemaType, _ := c.Flags().GetString("schema-type")
			customerID, _ := c.Flags().GetString("customer-id")
			content, err := parseContent(contentStr)
			if err != nil {
				return fmt.Errorf("invalid content: %w", err)
			}
			req := &client.CreateSchemaRequest{
				Name:        name,
				DocTypeCode: docType,
				Content:     content,
			}
			if description != "" {
				req.Description = &description
			}
			if country != "" {
				req.CountryCode = &country
			}
			if visibility != "" {
				vis := client.Visibility(visibility)
				req.Visibility = &vis
			}
			if schemaType != "" {
				st := client.SchemaType(schemaType)
				req.SchemaType = &st
			}
			if customerID != "" {
				req.CustomerID = &customerID
			}
			schema, err := GetClient().CreateSchema(req)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(schema)
			}
			output.PrintSuccess(fmt.Sprintf("Schema created: %s (version: %s)", schema.PublicID, schema.PublicVersionID))
			return nil
		},
	}
	createCmd.Flags().StringP("name", "n", "", "Schema name (required)")
	createCmd.Flags().StringP("doc-type", "t", "", "Doc type code (required)")
	createCmd.Flags().String("content", "", "JSON schema content or @filepath (required)")
	createCmd.Flags().StringP("description", "d", "", "Schema description")
	createCmd.Flags().StringP("country", "c", "", "Country code")
	createCmd.Flags().StringP("visibility", "v", "private", "Visibility (public, community, private)")
	createCmd.Flags().String("schema-type", "standard", "Schema type (standard, regex)")
	createCmd.Flags().String("customer-id", "", "Customer ID")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("doc-type")
	createCmd.MarkFlagRequired("content")

	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a schema",
		Long:  "Update a schema. If active, creates a new version.",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			req := &client.UpdateSchemaRequest{}
			hasUpdate := false
			if name, _ := c.Flags().GetString("name"); name != "" {
				req.Name = &name
				hasUpdate = true
			}
			if docType, _ := c.Flags().GetString("doc-type"); docType != "" {
				req.DocTypeCode = &docType
				hasUpdate = true
			}
			if contentStr, _ := c.Flags().GetString("content"); contentStr != "" {
				content, err := parseContent(contentStr)
				if err != nil {
					return fmt.Errorf("invalid content: %w", err)
				}
				req.Content = content
				hasUpdate = true
			}
			if description, _ := c.Flags().GetString("description"); c.Flags().Changed("description") {
				req.Description = &description
				hasUpdate = true
			}
			if country, _ := c.Flags().GetString("country"); c.Flags().Changed("country") {
				req.CountryCode = &country
				hasUpdate = true
			}
			if visibility, _ := c.Flags().GetString("visibility"); visibility != "" {
				vis := client.Visibility(visibility)
				req.Visibility = &vis
				hasUpdate = true
			}
			if schemaType, _ := c.Flags().GetString("schema-type"); schemaType != "" {
				st := client.SchemaType(schemaType)
				req.SchemaType = &st
				hasUpdate = true
			}
			if !hasUpdate {
				return fmt.Errorf("no update fields provided")
			}
			schema, err := GetClient().UpdateSchema(id, req)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(schema)
			}
			output.PrintSuccess(fmt.Sprintf("Schema updated: %s (version: %s)", schema.PublicID, schema.PublicVersionID))
			return nil
		},
	}
	updateCmd.Flags().StringP("name", "n", "", "Schema name")
	updateCmd.Flags().StringP("doc-type", "t", "", "Doc type code")
	updateCmd.Flags().String("content", "", "JSON schema content or @filepath")
	updateCmd.Flags().StringP("description", "d", "", "Schema description")
	updateCmd.Flags().StringP("country", "c", "", "Country code")
	updateCmd.Flags().StringP("visibility", "v", "", "Visibility (public, community, private)")
	updateCmd.Flags().String("schema-type", "", "Schema type (standard, regex)")

	activateCmd := &cobra.Command{
		Use:   "activate <id>",
		Short: "Activate a schema",
		Long:  "Transition a draft schema to active status",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			schema, err := GetClient().ActivateSchema(id)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(schema)
			}
			output.PrintSuccess(fmt.Sprintf("Schema activated: %s", schema.PublicID))
			return nil
		},
	}

	deprecateCmd := &cobra.Command{
		Use:   "deprecate <id>",
		Short: "Deprecate a schema",
		Long:  "Transition an active schema to deprecated status",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			schema, err := GetClient().DeprecateSchema(id)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(schema)
			}
			output.PrintSuccess(fmt.Sprintf("Schema deprecated: %s", schema.PublicID))
			return nil
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a schema",
		Long:  "Delete a draft schema. Active schemas must be deprecated first.",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			if err := GetClient().DeleteSchema(id); err != nil {
				return err
			}
			output.PrintSuccess(fmt.Sprintf("Schema deleted: %s", id))
			return nil
		},
	}

	findBestCmd := &cobra.Command{
		Use:   "find-best",
		Short: "Find best matching schema",
		Long:  "Find the best schema matching a doc type and optional country",
		RunE: func(c *cobra.Command, args []string) error {
			docType, _ := c.Flags().GetString("doc-type")
			country, _ := c.Flags().GetString("country")
			customerID, _ := c.Flags().GetString("customer-id")
			req := &client.FindBestRequest{
				DocTypeCode: docType,
			}
			if country != "" {
				req.CountryCode = &country
			}
			if customerID != "" {
				req.CustomerID = &customerID
			}
			result, err := GetClient().FindBestSchema(req)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(result)
			}
			if result.Schema == nil {
				output.PrintSuccess("No matching schema found")
				return nil
			}
			fmt.Printf("Match type: %s\n\n", output.PtrString(result.MatchType, "unknown"))
			printSchemaDetails(result.Schema)
			return nil
		},
	}
	findBestCmd.Flags().StringP("doc-type", "t", "", "Doc type code (required)")
	findBestCmd.Flags().StringP("country", "c", "", "Country code")
	findBestCmd.Flags().String("customer-id", "", "Customer ID")
	findBestCmd.MarkFlagRequired("doc-type")

	versionsCmd := &cobra.Command{
		Use:   "versions <id>",
		Short: "List schema versions",
		Long:  "List all versions of a schema by publicId",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			versions, err := GetClient().GetSchemaVersions(id)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(versions)
			}
			headers := []string{"VERSION ID", "VERSION", "STATUS", "CREATED AT"}
			rows := make([][]string, len(versions))
			for i, v := range versions {
				rows[i] = []string{
					v.PublicVersionID,
					strconv.Itoa(v.Version),
					string(v.Status),
					v.CreatedAt.Format("2006-01-02 15:04:05"),
				}
			}
			output.PrintTable(headers, rows)
			return nil
		},
	}

	matchCmd := &cobra.Command{
		Use:   "match <file>",
		Short: "Match schema to file",
		Long:  "Upload a PDF or JPEG file to classify and find matching schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			filePath := args[0]
			customerID, _ := c.Flags().GetString("customer-id")
			var custIDPtr *string
			if customerID != "" {
				custIDPtr = &customerID
			}
			result, err := GetClient().MatchSchema(filePath, custIDPtr)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(result)
			}
			fmt.Println("Classification:")
			fmt.Printf("  Doc Type: %s\n", result.Classification.DocType)
			fmt.Printf("  Country:  %s\n", result.Classification.Country)
			fmt.Printf("  Pages:    %v\n", result.Classification.Pages)
			fmt.Println()
			if result.Schema != nil {
				fmt.Println("Matched Schema:")
				fmt.Printf("  ID:         %s\n", result.Schema.PublicID)
				fmt.Printf("  Version ID: %s\n", result.Schema.PublicVersionID)
				fmt.Printf("  Name:       %s\n", result.Schema.Name)
				fmt.Printf("  Type:       %s\n", result.Schema.SchemaType)
			} else {
				fmt.Println("No matching schema found")
			}
			return nil
		},
	}
	matchCmd.Flags().String("customer-id", "", "Customer ID for private schema matching")

	cmd.AddCommand(listCmd, getCmd, createCmd, updateCmd, activateCmd, deprecateCmd, deleteCmd, findBestCmd, versionsCmd, matchCmd)
	return cmd
}

func buildDocTypesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "doc-types",
		Aliases: []string{"doctypes"},
		Short:   "Manage document types",
		Long:    "Commands for managing document types in the registry",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List document types",
		Long:  "List all document types. Use --all to include inactive ones.",
		RunE: func(c *cobra.Command, args []string) error {
			includeAll, _ := c.Flags().GetBool("all")
			docTypes, err := GetClient().ListDocTypes(includeAll)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(docTypes)
			}
			headers := []string{"CODE", "NAME", "DESCRIPTION", "ACTIVE"}
			rows := make([][]string, len(docTypes))
			for i, dt := range docTypes {
				rows[i] = []string{
					dt.Code,
					dt.Name,
					output.Truncate(output.PtrString(dt.Description, "-"), 40),
					output.BoolString(dt.IsActive),
				}
			}
			output.PrintTable(headers, rows)
			return nil
		},
	}
	listCmd.Flags().Bool("all", false, "Include inactive doc types")

	getCmd := &cobra.Command{
		Use:   "get <code>",
		Short: "Get a document type",
		Long:  "Get a document type by its code",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			code := args[0]
			docType, err := GetClient().GetDocType(code)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(docType)
			}
			fmt.Printf("Code:        %s\n", docType.Code)
			fmt.Printf("Name:        %s\n", docType.Name)
			fmt.Printf("Description: %s\n", output.PtrString(docType.Description, "-"))
			fmt.Printf("Active:      %s\n", output.BoolString(docType.IsActive))
			fmt.Printf("Created At:  %s\n", docType.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("Updated At:  %s\n", docType.UpdatedAt.Format("2006-01-02 15:04:05"))
			return nil
		},
	}

	createCmd := &cobra.Command{
		Use:   "create <code> <name>",
		Short: "Create a document type",
		Long:  "Create a new document type with the given code and name",
		Args:  cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			code := args[0]
			name := args[1]
			description, _ := c.Flags().GetString("description")
			req := &client.CreateDocTypeRequest{
				Code: code,
				Name: name,
			}
			if description != "" {
				req.Description = &description
			}
			docType, err := GetClient().CreateDocType(req)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(docType)
			}
			output.PrintSuccess(fmt.Sprintf("Doc type created: %s", docType.Code))
			return nil
		},
	}
	createCmd.Flags().StringP("description", "d", "", "Doc type description")

	updateCmd := &cobra.Command{
		Use:   "update <code>",
		Short: "Update a document type",
		Long:  "Update a document type's name, description, or active status",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			code := args[0]
			req := &client.UpdateDocTypeRequest{}
			hasUpdate := false
			if name, _ := c.Flags().GetString("name"); name != "" {
				req.Name = &name
				hasUpdate = true
			}
			if c.Flags().Changed("description") {
				description, _ := c.Flags().GetString("description")
				req.Description = &description
				hasUpdate = true
			}
			if c.Flags().Changed("active") {
				active, _ := c.Flags().GetBool("active")
				req.IsActive = &active
				hasUpdate = true
			}
			if !hasUpdate {
				return fmt.Errorf("no update fields provided")
			}
			docType, err := GetClient().UpdateDocType(code, req)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(docType)
			}
			output.PrintSuccess(fmt.Sprintf("Doc type updated: %s", docType.Code))
			return nil
		},
	}
	updateCmd.Flags().StringP("name", "n", "", "New name")
	updateCmd.Flags().StringP("description", "d", "", "New description")
	updateCmd.Flags().Bool("active", true, "Set active status")

	deleteCmd := &cobra.Command{
		Use:   "delete <code>",
		Short: "Delete a document type",
		Long:  "Soft delete a document type (sets isActive to false)",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			code := args[0]
			if err := GetClient().DeleteDocType(code); err != nil {
				return err
			}
			output.PrintSuccess(fmt.Sprintf("Doc type deleted: %s", code))
			return nil
		},
	}

	cmd.AddCommand(listCmd, getCmd, createCmd, updateCmd, deleteCmd)
	return cmd
}

func buildCountriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "countries",
		Short: "Manage countries",
		Long:  "Commands for managing countries in the registry",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List countries",
		Long:  "List all countries. Use --all to include inactive ones.",
		RunE: func(c *cobra.Command, args []string) error {
			includeAll, _ := c.Flags().GetBool("all")
			countries, err := GetClient().ListCountries(includeAll)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(countries)
			}
			headers := []string{"CODE", "NAME", "ACTIVE"}
			rows := make([][]string, len(countries))
			for i, country := range countries {
				rows[i] = []string{
					country.Code,
					country.Name,
					output.BoolString(country.IsActive),
				}
			}
			output.PrintTable(headers, rows)
			return nil
		},
	}
	listCmd.Flags().Bool("all", false, "Include inactive countries")

	getCmd := &cobra.Command{
		Use:   "get <code>",
		Short: "Get a country",
		Long:  "Get a country by its ISO 3166-1 alpha-2 code",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			code := args[0]
			country, err := GetClient().GetCountry(code)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(country)
			}
			fmt.Printf("Code:       %s\n", country.Code)
			fmt.Printf("Name:       %s\n", country.Name)
			fmt.Printf("Active:     %s\n", output.BoolString(country.IsActive))
			fmt.Printf("Created At: %s\n", country.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("Updated At: %s\n", country.UpdatedAt.Format("2006-01-02 15:04:05"))
			return nil
		},
	}

	createCmd := &cobra.Command{
		Use:   "create <code> <name>",
		Short: "Create a country",
		Long:  "Create a new country with the given ISO 3166-1 alpha-2 code and name",
		Args:  cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			code := args[0]
			name := args[1]
			req := &client.CreateCountryRequest{
				Code: code,
				Name: name,
			}
			country, err := GetClient().CreateCountry(req)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(country)
			}
			output.PrintSuccess(fmt.Sprintf("Country created: %s", country.Code))
			return nil
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update <code>",
		Short: "Update a country",
		Long:  "Update a country's name or active status",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			code := args[0]
			req := &client.UpdateCountryRequest{}
			hasUpdate := false
			if name, _ := c.Flags().GetString("name"); name != "" {
				req.Name = &name
				hasUpdate = true
			}
			if c.Flags().Changed("active") {
				active, _ := c.Flags().GetBool("active")
				req.IsActive = &active
				hasUpdate = true
			}
			if !hasUpdate {
				return fmt.Errorf("no update fields provided")
			}
			country, err := GetClient().UpdateCountry(code, req)
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(country)
			}
			output.PrintSuccess(fmt.Sprintf("Country updated: %s", country.Code))
			return nil
		},
	}
	updateCmd.Flags().StringP("name", "n", "", "New name")
	updateCmd.Flags().Bool("active", true, "Set active status")

	deleteCmd := &cobra.Command{
		Use:   "delete <code>",
		Short: "Delete a country",
		Long:  "Soft delete a country (sets isActive to false)",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			code := args[0]
			if err := GetClient().DeleteCountry(code); err != nil {
				return err
			}
			output.PrintSuccess(fmt.Sprintf("Country deleted: %s", code))
			return nil
		},
	}

	cmd.AddCommand(listCmd, getCmd, createCmd, updateCmd, deleteCmd)
	return cmd
}

func buildReferenceDataCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "reference-data",
		Aliases: []string{"ref", "reference"},
		Short:   "Get all reference data",
		Long:    "Get all active doc types and countries in a single request",
		RunE: func(c *cobra.Command, args []string) error {
			data, err := GetClient().GetReferenceData()
			if err != nil {
				return err
			}
			if output.JSONOutput {
				return output.PrintJSON(data)
			}
			fmt.Println("Document Types:")
			fmt.Println("---------------")
			headers := []string{"CODE", "NAME", "DESCRIPTION"}
			rows := make([][]string, len(data.DocTypes))
			for i, dt := range data.DocTypes {
				rows[i] = []string{
					dt.Code,
					dt.Name,
					output.Truncate(output.PtrString(dt.Description, "-"), 40),
				}
			}
			output.PrintTable(headers, rows)
			fmt.Println()
			fmt.Println("Countries:")
			fmt.Println("----------")
			headers = []string{"CODE", "NAME"}
			rows = make([][]string, len(data.Countries))
			for i, c := range data.Countries {
				rows[i] = []string{
					c.Code,
					c.Name,
				}
			}
			output.PrintTable(headers, rows)
			return nil
		},
	}
}

// initClient initializes the API client (called in shell mode)
func initClient() error {
	var err error
	cfg, err = config.Load()
	if err != nil {
		return err
	}
	cfg.WithAPIURL(apiURL).WithAPIKey(apiKey)
	if err := cfg.Validate(); err != nil {
		return err
	}
	output.JSONOutput = jsonOut
	apiClient = client.New(cfg)
	return nil
}

// printBanner prints the welcome ASCII art banner
func printBanner() {
	banner := `
  ███████╗ ██████╗██╗  ██╗███████╗███╗   ███╗ █████╗
  ██╔════╝██╔════╝██║  ██║██╔════╝████╗ ████║██╔══██╗
  ███████╗██║     ███████║█████╗  ██╔████╔██║███████║
  ╚════██║██║     ██╔══██║██╔══╝  ██║╚██╔╝██║██╔══██║
  ███████║╚██████╗██║  ██║███████╗██║ ╚═╝ ██║██║  ██║
  ╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝
  ██████╗ ███████╗ ██████╗ ██╗███████╗████████╗██████╗ ██╗   ██╗
  ██╔══██╗██╔════╝██╔════╝ ██║██╔════╝╚══██╔══╝██╔══██╗╚██╗ ██╔╝
  ██████╔╝█████╗  ██║  ███╗██║███████╗   ██║   ██████╔╝ ╚████╔╝
  ██╔══██╗██╔══╝  ██║   ██║██║╚════██║   ██║   ██╔══██╗  ╚██╔╝
  ██║  ██║███████╗╚██████╔╝██║███████║   ██║   ██║  ██║   ██║
  ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝   ╚═╝
`
	fmt.Println(banner)
	fmt.Printf("  Version: %s\n", Version)
	fmt.Println("  Type 'help' for commands, 'exit' to quit")
}

// printShellHelp prints the help message for the interactive shell
func printShellHelp() {
	help := `
COMMANDS
────────────────────────────────────────────────────────────────────────────────

  health              Check API health status
                      Example: health

  reference-data      Get all doc types and countries
                      Example: reference-data

SCHEMAS
────────────────────────────────────────────────────────────────────────────────

  schemas list        List all schemas (with optional filters)
                      Example: schemas list
                      Example: schemas list --status active
                      Example: schemas list --doc-type Invoice --country PT

  schemas get         Get a schema by ID
                      Example: schemas get sch_abc123def456

  schemas create      Create a new schema
                      Example: schemas create -n "My Schema" -t Invoice --content @schema.json

  schemas update      Update a schema
                      Example: schemas update sch_abc123def456 -n "New Name"

  schemas activate    Activate a draft schema
                      Example: schemas activate sch_abc123def456

  schemas deprecate   Deprecate an active schema
                      Example: schemas deprecate sch_abc123def456

  schemas delete      Delete a draft schema
                      Example: schemas delete sch_abc123def456

  schemas find-best   Find the best matching schema
                      Example: schemas find-best -t Invoice -c PT

  schemas versions    List all versions of a schema
                      Example: schemas versions sch_abc123def456

  schemas match       Upload a file to classify and match schema
                      Example: schemas match invoice.pdf

DOC-TYPES
────────────────────────────────────────────────────────────────────────────────

  doc-types list      List all document types
                      Example: doc-types list
                      Example: doc-types list --all    (include inactive)

  doc-types get       Get a document type by code
                      Example: doc-types get Invoice

  doc-types create    Create a new document type
                      Example: doc-types create Invoice "Invoice Document"

  doc-types update    Update a document type
                      Example: doc-types update Invoice -n "New Name"

  doc-types delete    Delete (deactivate) a document type
                      Example: doc-types delete Invoice

COUNTRIES
────────────────────────────────────────────────────────────────────────────────

  countries list      List all countries
                      Example: countries list
                      Example: countries list --all    (include inactive)

  countries get       Get a country by code
                      Example: countries get PT

  countries create    Create a new country
                      Example: countries create PT "Portugal"

  countries update    Update a country
                      Example: countries update PT -n "Portuguese Republic"

  countries delete    Delete (deactivate) a country
                      Example: countries delete PT

OTHER
────────────────────────────────────────────────────────────────────────────────

  version             Show version number
  help, ?             Show this help message
  exit, quit, q       Exit the shell

FLAGS (can be added to any command)
────────────────────────────────────────────────────────────────────────────────

  --json              Output results as JSON
                      Example: schemas list --json
`
	fmt.Println(help)
}
