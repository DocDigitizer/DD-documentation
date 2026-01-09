package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/miguel-bandeira-infosistema/schemactl/internal/client"
	"github.com/miguel-bandeira-infosistema/schemactl/internal/output"
	"github.com/spf13/cobra"
)

var schemasCmd = &cobra.Command{
	Use:   "schemas",
	Short: "Manage schemas",
	Long:  "Commands for managing JSON schemas in the registry",
}

func init() {
	schemasCmd.AddCommand(schemasListCmd)
	schemasCmd.AddCommand(schemasGetCmd)
	schemasCmd.AddCommand(schemasCreateCmd)
	schemasCmd.AddCommand(schemasUpdateCmd)
	schemasCmd.AddCommand(schemasActivateCmd)
	schemasCmd.AddCommand(schemasDeprecateCmd)
	schemasCmd.AddCommand(schemasDeleteCmd)
	schemasCmd.AddCommand(schemasFindBestCmd)
	schemasCmd.AddCommand(schemasVersionsCmd)
	schemasCmd.AddCommand(schemasMatchCmd)
	schemasCmd.AddCommand(schemasGenerateCmd)

	// List flags
	schemasListCmd.Flags().String("status", "", "Filter by status (draft, active, deprecated)")
	schemasListCmd.Flags().StringP("doc-type", "t", "", "Filter by doc type code")
	schemasListCmd.Flags().StringP("country", "c", "", "Filter by country code")
	schemasListCmd.Flags().StringP("visibility", "v", "", "Filter by visibility (public, community, private)")
	schemasListCmd.Flags().String("customer-id", "", "Filter by customer ID")
	schemasListCmd.Flags().Int("limit", 50, "Number of items to return")
	schemasListCmd.Flags().Int("offset", 0, "Number of items to skip")

	// Create flags
	schemasCreateCmd.Flags().StringP("name", "n", "", "Schema name (required)")
	schemasCreateCmd.Flags().StringP("doc-type", "t", "", "Doc type code (required)")
	schemasCreateCmd.Flags().String("content", "", "JSON schema content or @filepath (required)")
	schemasCreateCmd.Flags().StringP("description", "d", "", "Schema description")
	schemasCreateCmd.Flags().StringP("country", "c", "", "Country code")
	schemasCreateCmd.Flags().StringP("visibility", "v", "private", "Visibility (public, community, private)")
	schemasCreateCmd.Flags().String("schema-type", "standard", "Schema type (standard, regex)")
	schemasCreateCmd.Flags().String("customer-id", "", "Customer ID")
	schemasCreateCmd.MarkFlagRequired("name")
	schemasCreateCmd.MarkFlagRequired("doc-type")
	schemasCreateCmd.MarkFlagRequired("content")

	// Update flags
	schemasUpdateCmd.Flags().StringP("name", "n", "", "Schema name")
	schemasUpdateCmd.Flags().StringP("doc-type", "t", "", "Doc type code")
	schemasUpdateCmd.Flags().String("content", "", "JSON schema content or @filepath")
	schemasUpdateCmd.Flags().StringP("description", "d", "", "Schema description")
	schemasUpdateCmd.Flags().StringP("country", "c", "", "Country code")
	schemasUpdateCmd.Flags().StringP("visibility", "v", "", "Visibility (public, community, private)")
	schemasUpdateCmd.Flags().String("schema-type", "", "Schema type (standard, regex)")

	// Find-best flags
	schemasFindBestCmd.Flags().StringP("doc-type", "t", "", "Doc type code (required)")
	schemasFindBestCmd.Flags().StringP("country", "c", "", "Country code")
	schemasFindBestCmd.Flags().String("customer-id", "", "Customer ID")
	schemasFindBestCmd.MarkFlagRequired("doc-type")

	// Match flags
	schemasMatchCmd.Flags().String("customer-id", "", "Customer ID for private schema matching")

	// Generate flags
	schemasGenerateCmd.Flags().StringP("file", "f", "", "Path to PDF or JPEG file")
	schemasGenerateCmd.Flags().String("text", "", "Raw text content (alternative to file)")
	schemasGenerateCmd.Flags().StringP("doc-type", "t", "", "Doc type code (required)")
	schemasGenerateCmd.Flags().StringP("country", "c", "", "Country code (required)")
	schemasGenerateCmd.Flags().Bool("use-ocr", true, "Whether to use OCR on file (default: true)")
	schemasGenerateCmd.MarkFlagRequired("doc-type")
	schemasGenerateCmd.MarkFlagRequired("country")
}

var schemasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List schemas",
	Long:  "List schemas with optional filtering",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := &client.ListSchemasOptions{}

		if s, _ := cmd.Flags().GetString("status"); s != "" {
			status := client.Status(s)
			opts.Status = &status
		}
		if s, _ := cmd.Flags().GetString("doc-type"); s != "" {
			opts.DocType = &s
		}
		if s, _ := cmd.Flags().GetString("country"); s != "" {
			opts.Country = &s
		}
		if s, _ := cmd.Flags().GetString("visibility"); s != "" {
			vis := client.Visibility(s)
			opts.Visibility = &vis
		}
		if s, _ := cmd.Flags().GetString("customer-id"); s != "" {
			opts.CustomerID = &s
		}
		opts.Limit, _ = cmd.Flags().GetInt("limit")
		opts.Offset, _ = cmd.Flags().GetInt("offset")

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

var schemasGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a schema",
	Long:  "Get a schema by publicId (sch_xxx) or publicVersionId (schv_xxx)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

var schemasCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a schema",
	Long:  "Create a new schema in draft status",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		docType, _ := cmd.Flags().GetString("doc-type")
		contentStr, _ := cmd.Flags().GetString("content")
		description, _ := cmd.Flags().GetString("description")
		country, _ := cmd.Flags().GetString("country")
		visibility, _ := cmd.Flags().GetString("visibility")
		schemaType, _ := cmd.Flags().GetString("schema-type")
		customerID, _ := cmd.Flags().GetString("customer-id")

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

var schemasUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a schema",
	Long:  "Update a schema. If active, creates a new version.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		req := &client.UpdateSchemaRequest{}
		hasUpdate := false

		if name, _ := cmd.Flags().GetString("name"); name != "" {
			req.Name = &name
			hasUpdate = true
		}
		if docType, _ := cmd.Flags().GetString("doc-type"); docType != "" {
			req.DocTypeCode = &docType
			hasUpdate = true
		}
		if contentStr, _ := cmd.Flags().GetString("content"); contentStr != "" {
			content, err := parseContent(contentStr)
			if err != nil {
				return fmt.Errorf("invalid content: %w", err)
			}
			req.Content = content
			hasUpdate = true
		}
		if description, _ := cmd.Flags().GetString("description"); cmd.Flags().Changed("description") {
			req.Description = &description
			hasUpdate = true
		}
		if country, _ := cmd.Flags().GetString("country"); cmd.Flags().Changed("country") {
			req.CountryCode = &country
			hasUpdate = true
		}
		if visibility, _ := cmd.Flags().GetString("visibility"); visibility != "" {
			vis := client.Visibility(visibility)
			req.Visibility = &vis
			hasUpdate = true
		}
		if schemaType, _ := cmd.Flags().GetString("schema-type"); schemaType != "" {
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

var schemasActivateCmd = &cobra.Command{
	Use:   "activate <id>",
	Short: "Activate a schema",
	Long:  "Transition a draft schema to active status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

var schemasDeprecateCmd = &cobra.Command{
	Use:   "deprecate <id>",
	Short: "Deprecate a schema",
	Long:  "Transition an active schema to deprecated status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

var schemasDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a schema",
	Long:  "Delete a draft schema. Active schemas must be deprecated first.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		if err := GetClient().DeleteSchema(id); err != nil {
			return err
		}

		output.PrintSuccess(fmt.Sprintf("Schema deleted: %s", id))
		return nil
	},
}

var schemasFindBestCmd = &cobra.Command{
	Use:   "find-best",
	Short: "Find best matching schema",
	Long:  "Find the best schema matching a doc type and optional country",
	RunE: func(cmd *cobra.Command, args []string) error {
		docType, _ := cmd.Flags().GetString("doc-type")
		country, _ := cmd.Flags().GetString("country")
		customerID, _ := cmd.Flags().GetString("customer-id")

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

var schemasVersionsCmd = &cobra.Command{
	Use:   "versions <id>",
	Short: "List schema versions",
	Long:  "List all versions of a schema by publicId",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

var schemasMatchCmd = &cobra.Command{
	Use:   "match <file>",
	Short: "Match schema to file",
	Long:  "Upload a PDF or JPEG file to classify and find matching schema",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		customerID, _ := cmd.Flags().GetString("customer-id")

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

var schemasGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate schema from document",
	Long: `Generate a JSON schema from a document using LLM.

Use this when you know the document type and country, and need to generate
a schema for data extraction. Provide either a file or text content.

Examples:
  # Generate from a PDF file (with OCR)
  schemactl schemas generate -f invoice.pdf -t Invoice -c PT

  # Generate from an image without OCR (vision mode)
  schemactl schemas generate -f invoice.jpg -t Invoice -c US --use-ocr=false

  # Generate from text content
  schemactl schemas generate --text "Invoice Number: 12345" -t Invoice -c PT`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		text, _ := cmd.Flags().GetString("text")
		docType, _ := cmd.Flags().GetString("doc-type")
		country, _ := cmd.Flags().GetString("country")
		useOCR, _ := cmd.Flags().GetBool("use-ocr")

		if filePath == "" && text == "" {
			return fmt.Errorf("either --file or --text must be provided")
		}
		if filePath != "" && text != "" {
			return fmt.Errorf("only one of --file or --text can be provided")
		}

		req := &client.GenerateSchemaRequest{
			FilePath:    filePath,
			Text:        text,
			DocTypeCode: docType,
			CountryCode: country,
			UseOCR:      useOCR,
		}

		result, err := GetClient().GenerateSchema(req)
		if err != nil {
			return err
		}

		if output.JSONOutput {
			return output.PrintJSON(result)
		}

		fmt.Printf("Doc Type: %s\n", result.DocType)
		fmt.Printf("Country:  %s\n", result.Country)
		fmt.Println()
		fmt.Println("Generated Schema:")
		contentJSON, _ := json.MarshalIndent(result.Schema.Content, "  ", "  ")
		fmt.Printf("  %s\n", string(contentJSON))

		return nil
	},
}

// parseContent parses content from a string or file (prefixed with @)
func parseContent(s string) (map[string]interface{}, error) {
	var data []byte
	var err error

	if strings.HasPrefix(s, "@") {
		filePath := strings.TrimPrefix(s, "@")
		data, err = os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
	} else {
		data = []byte(s)
	}

	var content map[string]interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return content, nil
}

// printSchemaDetails prints schema details in a formatted way
func printSchemaDetails(s *client.SchemaWithRelations) {
	fmt.Printf("ID:           %s\n", s.PublicID)
	fmt.Printf("Version ID:   %s\n", s.PublicVersionID)
	fmt.Printf("Name:         %s\n", s.Name)
	fmt.Printf("Description:  %s\n", output.PtrString(s.Description, "-"))
	fmt.Printf("Version:      %d\n", s.Version)
	fmt.Printf("Status:       %s\n", s.Status)
	fmt.Printf("Doc Type:     %s (%s)\n", s.DocType.Code, s.DocType.Name)
	if s.Country != nil {
		fmt.Printf("Country:      %s (%s)\n", s.Country.Code, s.Country.Name)
	} else {
		fmt.Printf("Country:      -\n")
	}
	fmt.Printf("Visibility:   %s\n", s.Visibility)
	fmt.Printf("Schema Type:  %s\n", s.SchemaType)
	fmt.Printf("Customer ID:  %s\n", output.PtrString(s.CustomerID, "-"))
	fmt.Printf("Created At:   %s\n", s.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated At:   %s\n", s.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()
	fmt.Println("Content:")
	contentJSON, _ := json.MarshalIndent(s.Content, "  ", "  ")
	fmt.Printf("  %s\n", string(contentJSON))
}
