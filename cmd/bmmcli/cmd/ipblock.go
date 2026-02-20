// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/nvidia/bare-metal-manager-rest/client"
	"github.com/nvidia/bare-metal-manager-rest/cmd/bmmcli/internal/pagination"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ipBlockCmd = &cobra.Command{
	Use:   "ip-block",
	Short: "IP block operations",
}

var ipBlockListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all IP blocks",
	RunE:  runIPBlockList,
}

var ipBlockGetCmd = &cobra.Command{
	Use:   "get <ip-block-id>",
	Short: "Get IP block details",
	Args:  cobra.ExactArgs(1),
	RunE:  runIPBlockGet,
}

var ipBlockCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an IP block",
	RunE:  runIPBlockCreate,
}

var ipBlockUpdateCmd = &cobra.Command{
	Use:   "update <ip-block-id>",
	Short: "Update an IP block",
	Args:  cobra.ExactArgs(1),
	RunE:  runIPBlockUpdate,
}

var ipBlockDeleteCmd = &cobra.Command{
	Use:   "delete <ip-block-id>",
	Short: "Delete an IP block",
	Args:  cobra.ExactArgs(1),
	RunE:  runIPBlockDelete,
}

func init() {
	ipBlockListCmd.Flags().Bool("json", false, "output raw JSON")
	ipBlockListCmd.Flags().String("site-id", "", "filter by site ID")

	ipBlockCreateCmd.Flags().String("name", "", "name for the IP block (required)")
	ipBlockCreateCmd.Flags().String("site-id", "", "site ID (required)")
	ipBlockCreateCmd.Flags().String("prefix", "", "IP prefix e.g. 10.0.0.0 (required)")
	ipBlockCreateCmd.Flags().Int32("prefix-length", 0, "prefix length e.g. 16 (required)")
	ipBlockCreateCmd.MarkFlagRequired("name")
	ipBlockCreateCmd.MarkFlagRequired("site-id")
	ipBlockCreateCmd.MarkFlagRequired("prefix")
	ipBlockCreateCmd.MarkFlagRequired("prefix-length")

	ipBlockUpdateCmd.Flags().String("name", "", "new name")
	ipBlockUpdateCmd.Flags().String("description", "", "new description")

	rootCmd.AddCommand(ipBlockCmd)
	ipBlockCmd.AddCommand(ipBlockListCmd)
	ipBlockCmd.AddCommand(ipBlockGetCmd)
	ipBlockCmd.AddCommand(ipBlockCreateCmd)
	ipBlockCmd.AddCommand(ipBlockUpdateCmd)
	ipBlockCmd.AddCommand(ipBlockDeleteCmd)
}

func runIPBlockList(cmd *cobra.Command, args []string) error {
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

	blocks, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.IpBlock, *http.Response, error) {
		req := apiClient.IPBlockAPI.GetAllIpblock(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing IP blocks (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing IP blocks: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(blocks))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, blocks)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, blocks)
	default:
		return printIPBlockTable(os.Stdout, blocks)
	}
}

func runIPBlockGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	block, resp, err := apiClient.IPBlockAPI.GetIpblock(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting IP block (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting IP block: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, block)
	default:
		return printJSON(os.Stdout, block)
	}
}

func runIPBlockCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required")
	}

	name, _ := cmd.Flags().GetString("name")
	siteID, _ := cmd.Flags().GetString("site-id")
	prefix, _ := cmd.Flags().GetString("prefix")
	prefixLen, _ := cmd.Flags().GetInt32("prefix-length")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewIpBlockCreateRequest(name, siteID, "DatacenterOnly", prefix, prefixLen, "IPv4")

	block, resp, err := apiClient.IPBlockAPI.CreateIpblock(ctx, org).IpBlockCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating IP block (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating IP block: %v", err)
	}

	fmt.Fprintf(os.Stderr, "IP block created: %s (%s)\n", ptrStr(block.Name), ptrStr(block.Id))
	return printJSON(os.Stdout, block)
}

func runIPBlockUpdate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewIpBlockUpdateRequest()
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		req.SetName(name)
	}
	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		req.SetDescription(description)
	}

	block, resp, err := apiClient.IPBlockAPI.UpdateIpblock(ctx, org, args[0]).IpBlockUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating IP block (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating IP block: %v", err)
	}

	fmt.Fprintf(os.Stderr, "IP block updated: %s (%s)\n", ptrStr(block.Name), ptrStr(block.Id))
	return printJSON(os.Stdout, block)
}

func runIPBlockDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.IPBlockAPI.DeleteIpblock(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting IP block (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting IP block: %v", err)
	}

	fmt.Fprintf(os.Stderr, "IP block deleted: %s\n", args[0])
	return nil
}

func printIPBlockTable(w io.Writer, blocks []client.IpBlock) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE ID\tAGE\tID")

	for _, b := range blocks {
		name := ptrStr(b.Name)
		id := ptrStr(b.Id)
		siteID := ptrStr(b.SiteId)
		status := ""
		if b.Status != nil {
			status = string(*b.Status)
		}
		age := "<unknown>"
		if b.Created != nil {
			age = formatAge(time.Since(*b.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", name, status, siteID, age, id)
	}

	return tw.Flush()
}
