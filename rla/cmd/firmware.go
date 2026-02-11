package cmd

import (
	"github.com/spf13/cobra"
)

var firmwareCmd = &cobra.Command{
	Use:   "firmware",
	Short: "Firmware operations",
	Long:  `Commands for firmware management including version checking, updates, and scheduling.`,
}

func init() {
	rootCmd.AddCommand(firmwareCmd)
}
