package cmd

import (
	"github.com/spf13/cobra"
)

var powerCmd = &cobra.Command{
	Use:   "power",
	Short: "Power operations",
	Long:  `Commands for power management including control, status, and statistics.`,
}

func init() {
	rootCmd.AddCommand(powerCmd)
}
