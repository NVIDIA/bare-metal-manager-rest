package cmd

import (
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations",
	Long:  `Commands for database management such as migrations.`,
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
