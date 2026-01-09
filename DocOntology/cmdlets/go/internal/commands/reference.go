package commands

import (
	"fmt"

	"github.com/miguel-bandeira-infosistema/schemactl/internal/output"
	"github.com/spf13/cobra"
)

var referenceDataCmd = &cobra.Command{
	Use:     "reference-data",
	Aliases: []string{"ref", "reference"},
	Short:   "Get all reference data",
	Long:    "Get all active doc types and countries in a single request",
	RunE: func(cmd *cobra.Command, args []string) error {
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
