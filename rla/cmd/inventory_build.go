package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nvidia/carbide-rest/rla/internal/builder"
	"github.com/nvidia/carbide-rest/rla/pkg/client"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build the rack inventory",
		Long:  `Build the rack inventory from the files which are fetched from various sources`,

		Run: func(cmd *cobra.Command, args []string) {
			doBuild()
		},
	}

	sourceDumperDirs              []string
	buildInventoryDryRun          bool
	buildInventoryRLAHost         string
	buildInventoryRLAPort         int
	buildInventoryOutputDumperDir string
)

func init() {
	inventoryCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringArrayVarP(
		&sourceDumperDirs,
		"source",
		"s",
		[]string{},
		"source dumper directories",
	)

	buildCmd.Flags().BoolVarP(&buildInventoryDryRun, "dry-run", "d", false, "dry-run")                        //nolint
	buildCmd.Flags().StringVarP(&buildInventoryRLAHost, "host", "u", "localhost", "RLA service host")         //nolint
	buildCmd.Flags().IntVarP(&buildInventoryRLAPort, "port", "p", defaultServicePort, "RLA service port")     //nolint
	buildCmd.Flags().StringVarP(&buildInventoryOutputDumperDir, "output", "o", "", "output dumper directory") //nolint
}

func doBuild() {
	builder, err := builder.New(
		builder.Config{
			SourceDumperDirs: sourceDumperDirs,
			OutputDumperDir:  buildInventoryOutputDumperDir,
			DryRun:           buildInventoryDryRun,
			RLAClientConf: client.Config{
				Host: buildInventoryRLAHost,
				Port: buildInventoryRLAPort,
			},
		},
	)

	if err != nil {
		log.Fatal().Msgf("failed to build the builder: %v", err)
	}

	defer builder.Done()

	processed := builder.Build(context.Background())
	log.Info().Msgf("built %d racks, %d NVL domains", processed[0], processed[1])
}
