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

var subnetCmd = &cobra.Command{
	Use:   "subnet",
	Short: "Subnet operations",
}

var subnetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all subnets",
	RunE:  runSubnetList,
}

var subnetGetCmd = &cobra.Command{
	Use:   "get <subnet-id>",
	Short: "Get subnet details",
	Args:  cobra.ExactArgs(1),
	RunE:  runSubnetGet,
}

var subnetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a subnet",
	RunE:  runSubnetCreate,
}

var subnetUpdateCmd = &cobra.Command{
	Use:   "update <subnet-id>",
	Short: "Update a subnet",
	Args:  cobra.ExactArgs(1),
	RunE:  runSubnetUpdate,
}

var subnetDeleteCmd = &cobra.Command{
	Use:   "delete <subnet-id>",
	Short: "Delete a subnet",
	Args:  cobra.ExactArgs(1),
	RunE:  runSubnetDelete,
}

func init() {
	subnetListCmd.Flags().Bool("json", false, "output raw JSON")
	subnetListCmd.Flags().String("site-id", "", "filter by site ID")
	subnetListCmd.Flags().String("vpc-id", "", "filter by VPC ID")

	subnetCreateCmd.Flags().String("name", "", "name for the subnet (required)")
	subnetCreateCmd.Flags().String("vpc-id", "", "VPC ID (required)")
	subnetCreateCmd.Flags().String("ipv4-block-id", "", "IPv4 block ID (required)")
	subnetCreateCmd.Flags().Int32("prefix-length", 24, "prefix length (default 24)")
	subnetCreateCmd.MarkFlagRequired("name")
	subnetCreateCmd.MarkFlagRequired("vpc-id")
	subnetCreateCmd.MarkFlagRequired("ipv4-block-id")

	subnetUpdateCmd.Flags().String("name", "", "new name (required)")
	subnetUpdateCmd.Flags().String("description", "", "new description")
	subnetUpdateCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(subnetCmd)
	subnetCmd.AddCommand(subnetListCmd)
	subnetCmd.AddCommand(subnetGetCmd)
	subnetCmd.AddCommand(subnetCreateCmd)
	subnetCmd.AddCommand(subnetUpdateCmd)
	subnetCmd.AddCommand(subnetDeleteCmd)
}

func runSubnetList(cmd *cobra.Command, args []string) error {
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
	vpcID, _ := cmd.Flags().GetString("vpc-id")

	subnets, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.Subnet, *http.Response, error) {
		req := apiClient.SubnetAPI.GetAllSubnet(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		if vpcID != "" {
			req = req.VpcId(vpcID)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing subnets (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing subnets: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(subnets))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, subnets)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, subnets)
	default:
		return printSubnetTable(os.Stdout, subnets)
	}
}

func runSubnetGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	subnet, resp, err := apiClient.SubnetAPI.GetSubnet(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting subnet (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting subnet: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, subnet)
	default:
		return printJSON(os.Stdout, subnet)
	}
}

func runSubnetCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required")
	}

	name, _ := cmd.Flags().GetString("name")
	vpcID, _ := cmd.Flags().GetString("vpc-id")
	blockID, _ := cmd.Flags().GetString("ipv4-block-id")
	prefixLen, _ := cmd.Flags().GetInt32("prefix-length")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewSubnetCreateRequest(name, vpcID, prefixLen)
	req.Ipv4BlockId = &blockID

	subnet, resp, err := apiClient.SubnetAPI.CreateSubnet(ctx, org).SubnetCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating subnet (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating subnet: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Subnet created: %s (%s)\n", ptrStr(subnet.Name), ptrStr(subnet.Id))
	return printJSON(os.Stdout, subnet)
}

func runSubnetUpdate(cmd *cobra.Command, args []string) error {
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

	req := client.NewSubnetUpdateRequest(name)
	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		req.SetDescription(description)
	}

	subnet, resp, err := apiClient.SubnetAPI.UpdateSubnet(ctx, org, args[0]).SubnetUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating subnet (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating subnet: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Subnet updated: %s (%s)\n", ptrStr(subnet.Name), ptrStr(subnet.Id))
	return printJSON(os.Stdout, subnet)
}

func runSubnetDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.SubnetAPI.DeleteSubnet(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting subnet (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting subnet: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Subnet deleted: %s\n", args[0])
	return nil
}

func printSubnetTable(w io.Writer, subnets []client.Subnet) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tVPC ID\tAGE\tID")

	for _, s := range subnets {
		name := ptrStr(s.Name)
		id := ptrStr(s.Id)
		vpcID := ptrStr(s.VpcId)
		status := ""
		if s.Status != nil {
			status = string(*s.Status)
		}
		age := "<unknown>"
		if s.Created != nil {
			age = formatAge(time.Since(*s.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", name, status, vpcID, age, id)
	}

	return tw.Flush()
}
