package commands

import (
	"fmt"

	"github.com/miguel-bandeira-infosistema/schemactl/internal/output"
	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check API health",
	Long:  "Check the health status of the API server and database connection",
	RunE: func(cmd *cobra.Command, args []string) error {
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
