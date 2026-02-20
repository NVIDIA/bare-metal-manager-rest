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

var ibPartitionCmd = &cobra.Command{
	Use:   "infiniband-partition",
	Short: "InfiniBand partition operations",
}

var ibPartitionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all InfiniBand partitions",
	RunE:  runIBPartitionList,
}

var ibPartitionGetCmd = &cobra.Command{
	Use:   "get <partition-id>",
	Short: "Get InfiniBand partition details",
	Args:  cobra.ExactArgs(1),
	RunE:  runIBPartitionGet,
}

var ibPartitionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an InfiniBand partition",
	RunE:  runIBPartitionCreate,
}

var ibPartitionUpdateCmd = &cobra.Command{
	Use:   "update <partition-id>",
	Short: "Update an InfiniBand partition",
	Args:  cobra.ExactArgs(1),
	RunE:  runIBPartitionUpdate,
}

var ibPartitionDeleteCmd = &cobra.Command{
	Use:   "delete <partition-id>",
	Short: "Delete an InfiniBand partition",
	Args:  cobra.ExactArgs(1),
	RunE:  runIBPartitionDelete,
}

func init() {
	ibPartitionListCmd.Flags().Bool("json", false, "output raw JSON")
	ibPartitionListCmd.Flags().String("site-id", "", "filter by site ID")

	ibPartitionCreateCmd.Flags().String("name", "", "name for the partition (required)")
	ibPartitionCreateCmd.Flags().String("site-id", "", "site ID (required)")
	ibPartitionCreateCmd.Flags().String("description", "", "optional description")
	ibPartitionCreateCmd.MarkFlagRequired("name")
	ibPartitionCreateCmd.MarkFlagRequired("site-id")

	ibPartitionUpdateCmd.Flags().String("name", "", "new name (required)")
	ibPartitionUpdateCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(ibPartitionCmd)
	ibPartitionCmd.AddCommand(ibPartitionListCmd)
	ibPartitionCmd.AddCommand(ibPartitionGetCmd)
	ibPartitionCmd.AddCommand(ibPartitionCreateCmd)
	ibPartitionCmd.AddCommand(ibPartitionUpdateCmd)
	ibPartitionCmd.AddCommand(ibPartitionDeleteCmd)
}

func runIBPartitionList(cmd *cobra.Command, args []string) error {
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

	partitions, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.InfiniBandPartition, *http.Response, error) {
		req := apiClient.InfiniBandPartitionAPI.GetAllInfinibandPartition(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing InfiniBand partitions (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing InfiniBand partitions: %v", err)
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
		return printIBPartitionTable(os.Stdout, partitions)
	}
}

func runIBPartitionGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	partition, resp, err := apiClient.InfiniBandPartitionAPI.GetInfinibandPartition(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting InfiniBand partition (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting InfiniBand partition: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, partition)
	default:
		return printJSON(os.Stdout, partition)
	}
}

func runIBPartitionCreate(cmd *cobra.Command, args []string) error {
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

	req := client.NewInfiniBandPartitionCreateRequest(name, siteID)
	if description != "" {
		req.SetDescription(description)
	}

	partition, resp, err := apiClient.InfiniBandPartitionAPI.CreateInfinibandPartition(ctx, org).InfiniBandPartitionCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating InfiniBand partition (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating InfiniBand partition: %v", err)
	}

	fmt.Fprintf(os.Stderr, "InfiniBand partition created: %s (%s)\n", ptrStr(partition.Name), ptrStr(partition.Id))
	return printJSON(os.Stdout, partition)
}

func runIBPartitionUpdate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	name, _ := cmd.Flags().GetString("name")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewInfiniBandPartitionUpdateRequest(name)

	partition, resp, err := apiClient.InfiniBandPartitionAPI.UpdateInfinibandPartition(ctx, org, args[0]).InfiniBandPartitionUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating InfiniBand partition (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating InfiniBand partition: %v", err)
	}

	fmt.Fprintf(os.Stderr, "InfiniBand partition updated: %s (%s)\n", ptrStr(partition.Name), ptrStr(partition.Id))
	return printJSON(os.Stdout, partition)
}

func runIBPartitionDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.InfiniBandPartitionAPI.DeleteInfinibandPartition(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting InfiniBand partition (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting InfiniBand partition: %v", err)
	}

	fmt.Fprintf(os.Stderr, "InfiniBand partition deleted: %s\n", args[0])
	return nil
}

func printIBPartitionTable(w io.Writer, partitions []client.InfiniBandPartition) error {
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
