package cmd

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nvidia/carbide-rest/rla/internal/clients/cerebro"
	"github.com/nvidia/carbide-rest/rla/internal/dumper"
	cerebrofetcher "github.com/nvidia/carbide-rest/rla/internal/fetcher/cerebro"
	"github.com/nvidia/carbide-rest/rla/pkg/client"
)

var (
	cerebroCmd = &cobra.Command{
		Use:   "cerebro",
		Short: "Fetch the inventory from the Cerebro",
		Run: func(cmd *cobra.Command, args []string) {
			doFetchCerebro()
		},
	}

	cerebroRackNameFile string
	cerebroRackName     string
	cerebroURL          string
	cerebroAPIKey       string

	fetchCerebroDumperDir string
	fetchCerebroDryRun    bool
	fetchCerebroRLAHost   string
	fetchCerebroRLAPort   int
)

func init() {
	fetchCmd.AddCommand(cerebroCmd)

	cerebroCmd.Flags().StringVarP(&cerebroURL, "url", "u", "", "Cerebro URL")                                //nolint
	cerebroCmd.Flags().StringVarP(&cerebroAPIKey, "api-key", "k", "", "Cerebro API Key")                     //nolint
	cerebroCmd.Flags().StringVarP(&cerebroRackNameFile, "rack-name-file", "f", "", "Cerebro rack name file") //nolint
	cerebroCmd.Flags().StringVarP(&fetchCerebroDumperDir, "dumper-dir", "o", "", "Cerebro dumper directory") //nolint
	cerebroCmd.Flags().BoolVarP(&fetchCerebroDryRun, "dry-run", "d", false, "dry-run")                       //nolint
	cerebroCmd.Flags().StringVarP(&fetchCerebroRLAHost, "host", "s", "localhost", "RLA service host")        //nolint
	cerebroCmd.Flags().IntVarP(&fetchCerebroRLAPort, "port", "p", defaultServicePort, "RLA service port")    //nolint
	cerebroCmd.Flags().StringVarP(&cerebroRackName, "rack-name", "r", "", "Cerebro rack name")               //nolint
}

func doFetchCerebro() {
	cerebroFetcher, err := cerebrofetcher.New(
		context.Background(),
		cerebrofetcher.Config{
			RLAClientConf: client.Config{
				Host: fetchCerebroRLAHost,
				Port: fetchCerebroRLAPort,
			},
			CerebroClientConf: cerebro.Config{
				URL:    cerebroURL,
				APIKey: cerebroAPIKey,
			},
			DryRun:    fetchCerebroDryRun,
			DumperDir: fetchCerebroDumperDir,
		},
	)

	if err != nil {
		log.Fatal().Msgf("failed to build the fetcher: %v", err)
	}

	defer cerebroFetcher.Done()

	names := make([]string, 0)
	nameMap := make(map[string]struct{})

	if cerebroRackName != "" {
		names = append(names, cerebroRackName)
		nameMap[cerebroRackName] = struct{}{}
	}

	if cerebroRackNameFile != "" {
		data, err := os.ReadFile(cerebroRackNameFile)
		if err != nil {
			log.Fatal().Msgf("failed to read summary file: %v", err)
		}

		summary, err := dumper.NewSummary(data)
		if err != nil {
			log.Fatal().Msgf("failed to build summary: %v", err)
		}

		names = append(names, summary.RackNames...)
	}

	n := cerebroFetcher.Fetch(context.Background(), names)
	log.Info().Msgf("fetched %d racks out of %d", n, len(names))
}
