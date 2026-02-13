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

var instanceTypeCmd = &cobra.Command{
	Use:   "instance-type",
	Short: "Instance type operations",
}

var instanceTypeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all instance types",
	RunE:  runInstanceTypeList,
}

var instanceTypeGetCmd = &cobra.Command{
	Use:   "get <instance-type-id>",
	Short: "Get instance type details",
	Args:  cobra.ExactArgs(1),
	RunE:  runInstanceTypeGet,
}

var instanceTypeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an instance type",
	RunE:  runInstanceTypeCreate,
}

var instanceTypeUpdateCmd = &cobra.Command{
	Use:   "update <instance-type-id>",
	Short: "Update an instance type",
	Args:  cobra.ExactArgs(1),
	RunE:  runInstanceTypeUpdate,
}

var instanceTypeDeleteCmd = &cobra.Command{
	Use:   "delete <instance-type-id>",
	Short: "Delete an instance type",
	Args:  cobra.ExactArgs(1),
	RunE:  runInstanceTypeDelete,
}

func init() {
	instanceTypeListCmd.Flags().Bool("json", false, "output raw JSON")
	instanceTypeListCmd.Flags().String("site-id", "", "filter by site ID")
	instanceTypeListCmd.Flags().String("tenant-id", "", "filter by tenant ID")

	instanceTypeCreateCmd.Flags().String("name", "", "name for the instance type (required)")
	instanceTypeCreateCmd.Flags().String("site-id", "", "site ID (required)")
	instanceTypeCreateCmd.Flags().String("description", "", "optional description")
	instanceTypeCreateCmd.MarkFlagRequired("name")
	instanceTypeCreateCmd.MarkFlagRequired("site-id")

	instanceTypeUpdateCmd.Flags().String("name", "", "new name")
	instanceTypeUpdateCmd.Flags().String("description", "", "new description")

	rootCmd.AddCommand(instanceTypeCmd)
	instanceTypeCmd.AddCommand(instanceTypeListCmd)
	instanceTypeCmd.AddCommand(instanceTypeGetCmd)
	instanceTypeCmd.AddCommand(instanceTypeCreateCmd)
	instanceTypeCmd.AddCommand(instanceTypeUpdateCmd)
	instanceTypeCmd.AddCommand(instanceTypeDeleteCmd)
}

func runInstanceTypeList(cmd *cobra.Command, args []string) error {
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
	tenantID, _ := cmd.Flags().GetString("tenant-id")

	types, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.InstanceType, *http.Response, error) {
		req := apiClient.InstanceTypeAPI.GetAllInstanceType(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		if tenantID != "" {
			req = req.TenantId(tenantID)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing instance types (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing instance types: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(types))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, types)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, types)
	default:
		return printInstanceTypeTable(os.Stdout, types)
	}
}

func runInstanceTypeGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	it, resp, err := apiClient.InstanceTypeAPI.GetInstanceType(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting instance type (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting instance type: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, it)
	default:
		return printJSON(os.Stdout, it)
	}
}

func runInstanceTypeCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required")
	}

	name, _ := cmd.Flags().GetString("name")
	siteID, _ := cmd.Flags().GetString("site-id")
	description, _ := cmd.Flags().GetString("description")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewInstanceTypeCreateRequest(name, siteID)
	if description != "" {
		req.Description = &description
	}

	it, resp, err := apiClient.InstanceTypeAPI.CreateInstanceType(ctx, org).InstanceTypeCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating instance type (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating instance type: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Instance type created: %s (%s)\n", ptrStr(it.Name), ptrStr(it.Id))
	return printJSON(os.Stdout, it)
}

func runInstanceTypeUpdate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewInstanceTypeUpdateRequest()
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		req.SetName(name)
	}
	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		req.SetDescription(description)
	}

	it, resp, err := apiClient.InstanceTypeAPI.UpdateInstanceType(ctx, org, args[0]).InstanceTypeUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating instance type (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating instance type: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Instance type updated: %s (%s)\n", ptrStr(it.Name), ptrStr(it.Id))
	return printJSON(os.Stdout, it)
}

func runInstanceTypeDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.InstanceTypeAPI.DeleteInstanceType(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting instance type (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting instance type: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Instance type deleted: %s\n", args[0])
	return nil
}

func printInstanceTypeTable(w io.Writer, types []client.InstanceType) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE ID\tAGE\tID")

	for _, it := range types {
		name := ptrStr(it.Name)
		id := ptrStr(it.Id)
		siteID := ""
		if it.SiteId != nil {
			siteID = *it.SiteId
		}
		status := ""
		if it.Status != nil {
			status = string(*it.Status)
		}
		age := "<unknown>"
		if it.Created != nil {
			age = formatAge(time.Since(*it.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", name, status, siteID, age, id)
	}

	return tw.Flush()
}
