// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/nvidia/bare-metal-manager-rest/client"
	"github.com/nvidia/bare-metal-manager-rest/cmd/bmmcli/internal/pagination"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rackCmd = &cobra.Command{
	Use:   "rack",
	Short: "Rack operations",
}

var rackListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all racks",
	RunE:  runRackList,
}

var rackGetCmd = &cobra.Command{
	Use:   "get <rack-id>",
	Short: "Get rack details",
	Args:  cobra.ExactArgs(1),
	RunE:  runRackGet,
}

func init() {
	rackListCmd.Flags().Bool("json", false, "output raw JSON")
	rackListCmd.Flags().String("site-id", "", "filter by site ID")

	rootCmd.AddCommand(rackCmd)
	rackCmd.AddCommand(rackListCmd)
	rackCmd.AddCommand(rackGetCmd)
}

func runRackList(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}
	siteID, _ := cmd.Flags().GetString("site-id")

	racks, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.Rack, *http.Response, error) {
		req := apiClient.RackAPI.GetAllRack(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing racks (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing racks: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(racks))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, racks)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, racks)
	default:
		return printRackTable(os.Stdout, racks)
	}
}

func runRackGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	rack, resp, err := apiClient.RackAPI.GetRack(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting rack (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting rack: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, rack)
	default:
		return printJSON(os.Stdout, rack)
	}
}

func printRackTable(w io.Writer, racks []client.Rack) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tMANUFACTURER\tMODEL\tID")

	for _, r := range racks {
		name := ptrStr(r.Name)
		id := ptrStr(r.Id)
		manufacturer := ptrStr(r.Manufacturer)
		model := ptrStr(r.Model)

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", name, manufacturer, model, id)
	}

	return tw.Flush()
}
