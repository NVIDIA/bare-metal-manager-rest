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

var allocationCmd = &cobra.Command{
	Use:   "allocation",
	Short: "Allocation operations",
}

var allocationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all allocations",
	RunE:  runAllocationList,
}

var allocationGetCmd = &cobra.Command{
	Use:   "get <allocation-id>",
	Short: "Get allocation details",
	Args:  cobra.ExactArgs(1),
	RunE:  runAllocationGet,
}

var allocationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an allocation",
	RunE:  runAllocationCreate,
}

var allocationUpdateCmd = &cobra.Command{
	Use:   "update <allocation-id>",
	Short: "Update an allocation",
	Args:  cobra.ExactArgs(1),
	RunE:  runAllocationUpdate,
}

var allocationDeleteCmd = &cobra.Command{
	Use:   "delete <allocation-id>",
	Short: "Delete an allocation",
	Args:  cobra.ExactArgs(1),
	RunE:  runAllocationDelete,
}

func init() {
	allocationListCmd.Flags().Bool("json", false, "output raw JSON")
	allocationListCmd.Flags().String("site-id", "", "filter by site ID")

	allocationCreateCmd.Flags().String("name", "", "name for the allocation (required)")
	allocationCreateCmd.Flags().String("site-id", "", "site ID (required)")
	allocationCreateCmd.Flags().String("tenant-id", "", "tenant ID (required)")
	allocationCreateCmd.Flags().String("description", "", "optional description")
	allocationCreateCmd.MarkFlagRequired("name")
	allocationCreateCmd.MarkFlagRequired("site-id")
	allocationCreateCmd.MarkFlagRequired("tenant-id")

	allocationUpdateCmd.Flags().String("name", "", "new name")
	allocationUpdateCmd.Flags().String("description", "", "new description")

	rootCmd.AddCommand(allocationCmd)
	allocationCmd.AddCommand(allocationListCmd)
	allocationCmd.AddCommand(allocationGetCmd)
	allocationCmd.AddCommand(allocationCreateCmd)
	allocationCmd.AddCommand(allocationUpdateCmd)
	allocationCmd.AddCommand(allocationDeleteCmd)
}

func runAllocationList(cmd *cobra.Command, args []string) error {
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

	allocs, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.Allocation, *http.Response, error) {
		req := apiClient.AllocationAPI.GetAllAllocation(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing allocations (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing allocations: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(allocs))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, allocs)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, allocs)
	default:
		return printAllocationTable(os.Stdout, allocs)
	}
}

func runAllocationGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	alloc, resp, err := apiClient.AllocationAPI.GetAllocation(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting allocation (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting allocation: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, alloc)
	default:
		return printJSON(os.Stdout, alloc)
	}
}

func runAllocationCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	name, _ := cmd.Flags().GetString("name")
	siteID, _ := cmd.Flags().GetString("site-id")
	tenantID, _ := cmd.Flags().GetString("tenant-id")
	description, _ := cmd.Flags().GetString("description")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewAllocationCreateRequest(name, tenantID, siteID)
	if description != "" {
		req.SetDescription(description)
	}

	alloc, resp, err := apiClient.AllocationAPI.CreateAllocation(ctx, org).AllocationCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating allocation (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating allocation: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Allocation created: %s (%s)\n", ptrStr(alloc.Name), ptrStr(alloc.Id))
	return printJSON(os.Stdout, alloc)
}

func runAllocationUpdate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewAllocationUpdateRequest()
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		req.SetName(name)
	}
	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		req.SetDescription(description)
	}

	alloc, resp, err := apiClient.AllocationAPI.UpdateAllocation(ctx, org, args[0]).AllocationUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating allocation (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating allocation: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Allocation updated: %s (%s)\n", ptrStr(alloc.Name), ptrStr(alloc.Id))
	return printJSON(os.Stdout, alloc)
}

func runAllocationDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.AllocationAPI.DeleteAllocation(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting allocation (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting allocation: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Allocation deleted: %s\n", args[0])
	return nil
}

func printAllocationTable(w io.Writer, allocs []client.Allocation) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE ID\tAGE\tID")

	for _, a := range allocs {
		name := ptrStr(a.Name)
		id := ptrStr(a.Id)
		siteID := ptrStr(a.SiteId)
		status := ""
		if a.Status != nil {
			status = string(*a.Status)
		}
		age := "<unknown>"
		if a.Created != nil {
			age = formatAge(time.Since(*a.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", name, status, siteID, age, id)
	}

	return tw.Flush()
}
