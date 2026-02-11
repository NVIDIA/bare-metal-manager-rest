package cmd

import (
	"github.com/spf13/cobra"
)

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Inventory management",
	Long:  `Commands for managing rack inventory, including fetching from various sources and building.`,
}

func init() {
	rootCmd.AddCommand(inventoryCmd)
}
