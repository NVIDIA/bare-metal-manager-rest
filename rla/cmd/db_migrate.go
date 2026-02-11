package cmd

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nvidia/carbide-rest/rla/internal/db"
	"github.com/nvidia/carbide-rest/rla/internal/db/migrations"
	"github.com/nvidia/carbide-rest/rla/internal/db/postgres"
)

var (
	rollBack string

	// migrateCmd represents the migrate command
	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run the db migration",
		Long:  `Run the db migration`,
		Run: func(cmd *cobra.Command, args []string) {
			doMigration()
		},
	}
)

func init() {
	dbCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().StringVarP(&rollBack, "rollback", "r", "", "Roll back the schema to the way it was at the specified time.  This is the application time, not from the ID.  Format 2006-01-02T15:04:05")
}

func doMigration() {
	dbConf, err := db.BuildDBConfigFromEnv()
	if err != nil {
		log.Fatal().Msgf("Unable to build database configuration: %v", err)
	}

	ctx := context.Background()

	db, err := postgres.New(ctx, dbConf)
	if err != nil {
		log.Fatal().Msgf("failed to connect to DB: %v", err)
	}

	if rollBack != "" {
		rollbackTime, err := time.Parse("2006-01-02T15:04:05", rollBack)
		if err != nil {
			log.Fatal().Msg("Bad rollback time")
		}
		if err := migrations.Rollback(ctx, db, rollbackTime); err != nil {
			log.Fatal().Msgf("Failed to roll back migrations: %v", err)
		}
	} else {
		if err := migrations.Migrate(ctx, db); err != nil {
			log.Fatal().Msgf("Failed to run migrations: %v", err)
		}
	}
}
