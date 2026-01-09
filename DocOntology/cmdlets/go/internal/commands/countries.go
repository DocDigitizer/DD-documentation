package commands

import (
	"fmt"

	"github.com/miguel-bandeira-infosistema/schemactl/internal/client"
	"github.com/miguel-bandeira-infosistema/schemactl/internal/output"
	"github.com/spf13/cobra"
)

var countriesCmd = &cobra.Command{
	Use:   "countries",
	Short: "Manage countries",
	Long:  "Commands for managing countries in the registry",
}

func init() {
	countriesCmd.AddCommand(countriesListCmd)
	countriesCmd.AddCommand(countriesGetCmd)
	countriesCmd.AddCommand(countriesCreateCmd)
	countriesCmd.AddCommand(countriesUpdateCmd)
	countriesCmd.AddCommand(countriesDeleteCmd)

	// List flags
	countriesListCmd.Flags().Bool("all", false, "Include inactive countries")

	// Update flags
	countriesUpdateCmd.Flags().StringP("name", "n", "", "New name")
	countriesUpdateCmd.Flags().Bool("active", true, "Set active status")
}

var countriesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List countries",
	Long:  "List all countries. Use --all to include inactive ones.",
	RunE: func(cmd *cobra.Command, args []string) error {
		includeAll, _ := cmd.Flags().GetBool("all")

		countries, err := GetClient().ListCountries(includeAll)
		if err != nil {
			return err
		}

		if output.JSONOutput {
			return output.PrintJSON(countries)
		}

		headers := []string{"CODE", "NAME", "ACTIVE"}
		rows := make([][]string, len(countries))
		for i, c := range countries {
			rows[i] = []string{
				c.Code,
				c.Name,
				output.BoolString(c.IsActive),
			}
		}
		output.PrintTable(headers, rows)

		return nil
	},
}

var countriesGetCmd = &cobra.Command{
	Use:   "get <code>",
	Short: "Get a country",
	Long:  "Get a country by its ISO 3166-1 alpha-2 code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

var countriesCreateCmd = &cobra.Command{
	Use:   "create <code> <name>",
	Short: "Create a country",
	Long:  "Create a new country with the given ISO 3166-1 alpha-2 code and name",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
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

var countriesUpdateCmd = &cobra.Command{
	Use:   "update <code>",
	Short: "Update a country",
	Long:  "Update a country's name or active status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]

		req := &client.UpdateCountryRequest{}
		hasUpdate := false

		if name, _ := cmd.Flags().GetString("name"); name != "" {
			req.Name = &name
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

var countriesDeleteCmd = &cobra.Command{
	Use:   "delete <code>",
	Short: "Delete a country",
	Long:  "Soft delete a country (sets isActive to false)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]

		if err := GetClient().DeleteCountry(code); err != nil {
			return err
		}

		output.PrintSuccess(fmt.Sprintf("Country deleted: %s", code))
		return nil
	},
}
