package cmd

import (
	"github.com/spf13/cobra"
)

var componentCmd = &cobra.Command{
	Use:   "component",
	Short: "Component operations",
	Long:  `Commands for querying and comparing components (expected vs actual).`,
}

func init() {
	rootCmd.AddCommand(componentCmd)
}
