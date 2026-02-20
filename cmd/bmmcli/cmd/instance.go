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

var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Instance operations",
}

var instanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all instances",
	RunE:  runInstanceList,
}

var instanceGetCmd = &cobra.Command{
	Use:   "get <instance-id>",
	Short: "Get instance details",
	Args:  cobra.ExactArgs(1),
	RunE:  runInstanceGet,
}

var instanceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an instance",
	RunE:  runInstanceCreate,
}

var instanceUpdateCmd = &cobra.Command{
	Use:   "update <instance-id>",
	Short: "Update an instance",
	Args:  cobra.ExactArgs(1),
	RunE:  runInstanceUpdate,
}

var instanceDeleteCmd = &cobra.Command{
	Use:   "delete <instance-id>",
	Short: "Delete an instance",
	Args:  cobra.ExactArgs(1),
	RunE:  runInstanceDelete,
}

func init() {
	instanceListCmd.Flags().Bool("json", false, "output raw JSON")
	instanceListCmd.Flags().String("vpc-id", "", "filter by VPC ID")
	instanceListCmd.Flags().String("site-id", "", "filter by site ID")
	instanceListCmd.Flags().String("status", "", "filter by status")
	instanceListCmd.Flags().String("name", "", "filter by name")

	instanceCreateCmd.Flags().String("name", "", "name for the instance (required)")
	instanceCreateCmd.Flags().String("vpc-id", "", "VPC ID (required)")
	instanceCreateCmd.Flags().String("subnet-id", "", "subnet ID for the network interface (required)")
	instanceCreateCmd.Flags().String("instance-type-id", "", "instance type ID (required)")
	instanceCreateCmd.Flags().String("tenant-id", "", "tenant ID (required)")
	instanceCreateCmd.Flags().String("os-id", "", "operating system ID (optional)")
	instanceCreateCmd.Flags().String("ipxe-script", "", "iPXE script (optional, used if no OS)")
	instanceCreateCmd.MarkFlagRequired("name")
	instanceCreateCmd.MarkFlagRequired("vpc-id")
	instanceCreateCmd.MarkFlagRequired("subnet-id")
	instanceCreateCmd.MarkFlagRequired("instance-type-id")
	instanceCreateCmd.MarkFlagRequired("tenant-id")

	instanceUpdateCmd.Flags().String("name", "", "new name")
	instanceUpdateCmd.Flags().String("description", "", "new description")

	rootCmd.AddCommand(instanceCmd)
	instanceCmd.AddCommand(instanceListCmd)
	instanceCmd.AddCommand(instanceGetCmd)
	instanceCmd.AddCommand(instanceCreateCmd)
	instanceCmd.AddCommand(instanceUpdateCmd)
	instanceCmd.AddCommand(instanceDeleteCmd)
}

func runInstanceList(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	vpcID, _ := cmd.Flags().GetString("vpc-id")
	siteID, _ := cmd.Flags().GetString("site-id")
	status, _ := cmd.Flags().GetString("status")
	name, _ := cmd.Flags().GetString("name")

	instances, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.Instance, *http.Response, error) {
		req := apiClient.InstanceAPI.GetAllInstance(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if vpcID != "" {
			req = req.VpcId(vpcID)
		}
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		if status != "" {
			req = req.Status(status)
		}
		if name != "" {
			req = req.Name(name)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing instances (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing instances: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(instances))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, instances)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, instances)
	default:
		return printInstanceTable(os.Stdout, instances)
	}
}

func runInstanceGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	instance, resp, err := apiClient.InstanceAPI.GetInstance(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting instance (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting instance: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, instance)
	default:
		return printJSON(os.Stdout, instance)
	}
}

func runInstanceUpdate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewInstanceUpdateRequest()
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		req.SetName(name)
	}
	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		req.SetDescription(description)
	}

	instance, resp, err := apiClient.InstanceAPI.UpdateInstance(ctx, org, args[0]).InstanceUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating instance (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating instance: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Instance updated: %s (%s)\n", ptrStr(instance.Name), ptrStr(instance.Id))
	return printJSON(os.Stdout, instance)
}

func runInstanceDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.InstanceAPI.DeleteInstance(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting instance (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting instance: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Instance deleted: %s\n", args[0])
	return nil
}

func runInstanceCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required")
	}

	name, _ := cmd.Flags().GetString("name")
	vpcID, _ := cmd.Flags().GetString("vpc-id")
	subnetID, _ := cmd.Flags().GetString("subnet-id")
	itID, _ := cmd.Flags().GetString("instance-type-id")
	tenantID, _ := cmd.Flags().GetString("tenant-id")
	osID, _ := cmd.Flags().GetString("os-id")
	ipxeScript, _ := cmd.Flags().GetString("ipxe-script")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	iface := client.InterfaceCreateRequest{
		SubnetId: &subnetID,
	}

	req := client.NewInstanceCreateRequest(name, tenantID, vpcID, []client.InterfaceCreateRequest{iface})
	req.InstanceTypeId = &itID

	if osID != "" {
		req.SetOperatingSystemId(osID)
	} else if ipxeScript != "" {
		req.SetIpxeScript(ipxeScript)
	} else {
		req.SetIpxeScript("#!ipxe\necho No OS configured")
	}

	instance, resp, err := apiClient.InstanceAPI.CreateInstance(ctx, org).InstanceCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating instance (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating instance: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Instance created: %s (%s)\n", ptrStr(instance.Name), ptrStr(instance.Id))
	return printJSON(os.Stdout, instance)
}

func printInstanceTable(w io.Writer, instances []client.Instance) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE ID\tVPC ID\tAGE\tID")

	for _, inst := range instances {
		name := ptrStr(inst.Name)
		id := ptrStr(inst.Id)
		siteID := ptrStr(inst.SiteId)
		vpcID := ptrStr(inst.VpcId)
		status := ""
		if inst.Status != nil {
			status = string(*inst.Status)
		}
		age := "<unknown>"
		if inst.Created != nil {
			age = formatAge(time.Since(*inst.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", name, status, siteID, vpcID, age, id)
	}

	return tw.Flush()
}
