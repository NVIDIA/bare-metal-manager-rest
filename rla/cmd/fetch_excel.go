package cmd

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nvidia/carbide-rest/rla/internal/fetcher/excel"
	"github.com/nvidia/carbide-rest/rla/pkg/client"
)

var (
	excelCmd = &cobra.Command{
		Use:   "excel",
		Short: "Fetch the inventory from the excel files",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(rackFileAndSheet) != 2 {
				return errors.New("rack excel information is incorrect")
			}

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			doFetchExcel()
		},
	}

	rackFileAndSheet []string
	compFileAndSheet []string

	fetchExcelDryRun    bool
	fetchExcelRLAHost   string
	fetchExcelRLAPort   int
	fetchExcelDumperDir string
)

func init() {
	fetchCmd.AddCommand(excelCmd)

	excelCmd.Flags().StringArrayVarP(
		&rackFileAndSheet,
		"rack",
		"r",
		[]string{},
		"excel file and sheet for racks",
	)

	excelCmd.Flags().StringArrayVarP(
		&compFileAndSheet,
		"component",
		"c",
		[]string{},
		"excel file and sheet for components",
	)

	excelCmd.Flags().BoolVarP(&fetchExcelDryRun, "dry-run", "d", false, "dry-run")                    //nolint
	excelCmd.Flags().StringVarP(&fetchExcelRLAHost, "host", "s", "localhost", "RLA service host")     //nolint
	excelCmd.Flags().IntVarP(&fetchExcelRLAPort, "port", "p", defaultServicePort, "RLA service port") //nolint
	excelCmd.Flags().StringVarP(&fetchExcelDumperDir, "dumper-dir", "o", "", "dumper directory")      //nolint
}

func doFetchExcel() {
	buildExcelInfo := func(args []string) excel.Info {
		var info excel.Info
		if len(args) == 2 {
			info.File, info.Sheet = args[0], args[1]
		}

		return info
	}

	excelFetcher, err := excel.New(
		excel.Config{
			RackExcelInfo: buildExcelInfo(rackFileAndSheet),
			CompExcelInfo: buildExcelInfo(compFileAndSheet),
			DumperDir:     fetchExcelDumperDir,
			RLAClientConf: client.Config{
				Host: fetchExcelRLAHost,
				Port: fetchExcelRLAPort,
			},
			DryRun: fetchExcelDryRun,
		},
	)

	if err != nil {
		log.Fatal().Msgf("failed to build the fetcher: %v", err)
	}

	defer excelFetcher.Done()

	if err := excelFetcher.Fetch(context.Background()); err != nil {
		log.Fatal().Msgf("failed to fetch from excel: %v", err)
	}
}
