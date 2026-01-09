package main

import (
	"os"

	"github.com/miguel-bandeira-infosistema/schemactl/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
