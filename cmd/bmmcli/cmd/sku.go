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

var skuCmd = &cobra.Command{
	Use:   "sku",
	Short: "SKU operations",
}

var skuListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all SKUs",
	RunE:  runSkuList,
}

var skuGetCmd = &cobra.Command{
	Use:   "get <sku-id>",
	Short: "Get SKU details",
	Args:  cobra.ExactArgs(1),
	RunE:  runSkuGet,
}

func init() {
	skuListCmd.Flags().Bool("json", false, "output raw JSON")
	skuListCmd.Flags().String("site-id", "", "filter by site ID")

	rootCmd.AddCommand(skuCmd)
	skuCmd.AddCommand(skuListCmd)
	skuCmd.AddCommand(skuGetCmd)
}

func runSkuList(cmd *cobra.Command, args []string) error {
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

	skus, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.Sku, *http.Response, error) {
		req := apiClient.SKUAPI.GetAllSku(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing SKUs (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing SKUs: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(skus))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, skus)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, skus)
	default:
		return printSkuTable(os.Stdout, skus)
	}
}

func runSkuGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	sku, resp, err := apiClient.SKUAPI.GetSku(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting SKU (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting SKU: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, sku)
	default:
		return printJSON(os.Stdout, sku)
	}
}

func printSkuTable(w io.Writer, skus []client.Sku) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "SITE ID\tDEVICE TYPE\tID")

	for _, s := range skus {
		id := ptrStr(s.Id)
		siteID := ptrStr(s.SiteId)
		deviceType := ""
		if s.DeviceType.IsSet() {
			deviceType = *s.DeviceType.Get()
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\n", siteID, deviceType, id)
	}

	return tw.Flush()
}
