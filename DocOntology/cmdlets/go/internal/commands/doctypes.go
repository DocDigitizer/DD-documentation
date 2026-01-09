package commands

import (
	"fmt"

	"github.com/miguel-bandeira-infosistema/schemactl/internal/client"
	"github.com/miguel-bandeira-infosistema/schemactl/internal/output"
	"github.com/spf13/cobra"
)

var docTypesCmd = &cobra.Command{
	Use:     "doc-types",
	Aliases: []string{"doctypes"},
	Short:   "Manage document types",
	Long:    "Commands for managing document types in the registry",
}

func init() {
	docTypesCmd.AddCommand(docTypesListCmd)
	docTypesCmd.AddCommand(docTypesGetCmd)
	docTypesCmd.AddCommand(docTypesCreateCmd)
	docTypesCmd.AddCommand(docTypesUpdateCmd)
	docTypesCmd.AddCommand(docTypesDeleteCmd)

	// List flags
	docTypesListCmd.Flags().Bool("all", false, "Include inactive doc types")

	// Create flags
	docTypesCreateCmd.Flags().StringP("description", "d", "", "Doc type description")

	// Update flags
	docTypesUpdateCmd.Flags().StringP("name", "n", "", "New name")
	docTypesUpdateCmd.Flags().StringP("description", "d", "", "New description")
	docTypesUpdateCmd.Flags().Bool("active", true, "Set active status")
}

var docTypesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List document types",
	Long:  "List all document types. Use --all to include inactive ones.",
	RunE: func(cmd *cobra.Command, args []string) error {
		includeAll, _ := cmd.Flags().GetBool("all")

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

var docTypesGetCmd = &cobra.Command{
	Use:   "get <code>",
	Short: "Get a document type",
	Long:  "Get a document type by its code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

var docTypesCreateCmd = &cobra.Command{
	Use:   "create <code> <name>",
	Short: "Create a document type",
	Long:  "Create a new document type with the given code and name",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]
		name := args[1]
		description, _ := cmd.Flags().GetString("description")

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

var docTypesUpdateCmd = &cobra.Command{
	Use:   "update <code>",
	Short: "Update a document type",
	Long:  "Update a document type's name, description, or active status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]

		req := &client.UpdateDocTypeRequest{}
		hasUpdate := false

		if name, _ := cmd.Flags().GetString("name"); name != "" {
			req.Name = &name
			hasUpdate = true
		}
		if cmd.Flags().Changed("description") {
			description, _ := cmd.Flags().GetString("description")
			req.Description = &description
			hasUpdate = true
		}
		if cmd.Flags().Changed("active") {
			active, _ := cmd.Flags().GetBool("active")
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

var docTypesDeleteCmd = &cobra.Command{
	Use:   "delete <code>",
	Short: "Delete a document type",
	Long:  "Soft delete a document type (sets isActive to false)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]

		if err := GetClient().DeleteDocType(code); err != nil {
			return err
		}

		output.PrintSuccess(fmt.Sprintf("Doc type deleted: %s", code))
		return nil
	},
}
