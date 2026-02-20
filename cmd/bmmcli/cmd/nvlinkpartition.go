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

var nvlinkPartitionCmd = &cobra.Command{
	Use:   "nvlink-logical-partition",
	Short: "NVLink logical partition operations",
}

var nvlinkPartitionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all NVLink logical partitions",
	RunE:  runNVLinkPartitionList,
}

var nvlinkPartitionGetCmd = &cobra.Command{
	Use:   "get <partition-id>",
	Short: "Get NVLink logical partition details",
	Args:  cobra.ExactArgs(1),
	RunE:  runNVLinkPartitionGet,
}

var nvlinkPartitionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an NVLink logical partition",
	RunE:  runNVLinkPartitionCreate,
}

var nvlinkPartitionUpdateCmd = &cobra.Command{
	Use:   "update <partition-id>",
	Short: "Update an NVLink logical partition",
	Args:  cobra.ExactArgs(1),
	RunE:  runNVLinkPartitionUpdate,
}

var nvlinkPartitionDeleteCmd = &cobra.Command{
	Use:   "delete <partition-id>",
	Short: "Delete an NVLink logical partition",
	Args:  cobra.ExactArgs(1),
	RunE:  runNVLinkPartitionDelete,
}

func init() {
	nvlinkPartitionListCmd.Flags().Bool("json", false, "output raw JSON")
	nvlinkPartitionListCmd.Flags().String("site-id", "", "filter by site ID")

	nvlinkPartitionCreateCmd.Flags().String("name", "", "name for the partition (required)")
	nvlinkPartitionCreateCmd.Flags().String("site-id", "", "site ID (required)")
	nvlinkPartitionCreateCmd.Flags().String("description", "", "optional description")
	nvlinkPartitionCreateCmd.MarkFlagRequired("name")
	nvlinkPartitionCreateCmd.MarkFlagRequired("site-id")

	nvlinkPartitionUpdateCmd.Flags().String("name", "", "new name")
	nvlinkPartitionUpdateCmd.Flags().String("description", "", "new description")

	rootCmd.AddCommand(nvlinkPartitionCmd)
	nvlinkPartitionCmd.AddCommand(nvlinkPartitionListCmd)
	nvlinkPartitionCmd.AddCommand(nvlinkPartitionGetCmd)
	nvlinkPartitionCmd.AddCommand(nvlinkPartitionCreateCmd)
	nvlinkPartitionCmd.AddCommand(nvlinkPartitionUpdateCmd)
	nvlinkPartitionCmd.AddCommand(nvlinkPartitionDeleteCmd)
}

func runNVLinkPartitionList(cmd *cobra.Command, args []string) error {
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

	partitions, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.NVLinkLogicalPartition, *http.Response, error) {
		req := apiClient.NVLinkLogicalPartitionAPI.GetAllNvlinkLogicalPartition(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing NVLink logical partitions (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing NVLink logical partitions: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(partitions))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, partitions)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, partitions)
	default:
		return printNVLinkPartitionTable(os.Stdout, partitions)
	}
}

func runNVLinkPartitionGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	partition, resp, err := apiClient.NVLinkLogicalPartitionAPI.GetNvlinkLogicalPartition(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting NVLink logical partition (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting NVLink logical partition: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, partition)
	default:
		return printJSON(os.Stdout, partition)
	}
}

func runNVLinkPartitionCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	name, _ := cmd.Flags().GetString("name")
	siteID, _ := cmd.Flags().GetString("site-id")
	description, _ := cmd.Flags().GetString("description")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewNVLinkLogicalPartitionCreateRequest(name, siteID)
	if description != "" {
		req.SetDescription(description)
	}

	partition, resp, err := apiClient.NVLinkLogicalPartitionAPI.CreateNvlinkLogicalPartition(ctx, org).NVLinkLogicalPartitionCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating NVLink logical partition (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating NVLink logical partition: %v", err)
	}

	fmt.Fprintf(os.Stderr, "NVLink logical partition created: %s (%s)\n", ptrStr(partition.Name), ptrStr(partition.Id))
	return printJSON(os.Stdout, partition)
}

func runNVLinkPartitionUpdate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewNVLinkLogicalPartitionUpdateRequest()
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		req.SetName(name)
	}
	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		req.SetDescription(description)
	}

	partition, resp, err := apiClient.NVLinkLogicalPartitionAPI.UpdateNvlinkLogicalPartition(ctx, org, args[0]).NVLinkLogicalPartitionUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating NVLink logical partition (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating NVLink logical partition: %v", err)
	}

	fmt.Fprintf(os.Stderr, "NVLink logical partition updated: %s (%s)\n", ptrStr(partition.Name), ptrStr(partition.Id))
	return printJSON(os.Stdout, partition)
}

func runNVLinkPartitionDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.NVLinkLogicalPartitionAPI.DeleteNvlinkLogicalPartition(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting NVLink logical partition (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting NVLink logical partition: %v", err)
	}

	fmt.Fprintf(os.Stderr, "NVLink logical partition deleted: %s\n", args[0])
	return nil
}

func printNVLinkPartitionTable(w io.Writer, partitions []client.NVLinkLogicalPartition) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE ID\tAGE\tID")

	for _, p := range partitions {
		name := ptrStr(p.Name)
		id := ptrStr(p.Id)
		siteID := ptrStr(p.SiteId)
		status := ""
		if p.Status != nil {
			status = string(*p.Status)
		}
		age := "<unknown>"
		if p.Created != nil {
			age = formatAge(time.Since(*p.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", name, status, siteID, age, id)
	}

	return tw.Flush()
}
