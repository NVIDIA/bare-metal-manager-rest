package cmd

import (
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch inventory from various sources",
	Long:  `Commands for fetching inventory data from sources like Cerebro or Excel files.`,
}

func init() {
	inventoryCmd.AddCommand(fetchCmd)
}
