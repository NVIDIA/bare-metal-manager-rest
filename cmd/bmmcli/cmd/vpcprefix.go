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

var vpcPrefixCmd = &cobra.Command{
	Use:   "vpc-prefix",
	Short: "VPC prefix operations",
}

var vpcPrefixListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VPC prefixes",
	RunE:  runVpcPrefixList,
}

var vpcPrefixGetCmd = &cobra.Command{
	Use:   "get <vpc-prefix-id>",
	Short: "Get VPC prefix details",
	Args:  cobra.ExactArgs(1),
	RunE:  runVpcPrefixGet,
}

var vpcPrefixCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a VPC prefix",
	RunE:  runVpcPrefixCreate,
}

var vpcPrefixUpdateCmd = &cobra.Command{
	Use:   "update <vpc-prefix-id>",
	Short: "Update a VPC prefix",
	Args:  cobra.ExactArgs(1),
	RunE:  runVpcPrefixUpdate,
}

var vpcPrefixDeleteCmd = &cobra.Command{
	Use:   "delete <vpc-prefix-id>",
	Short: "Delete a VPC prefix",
	Args:  cobra.ExactArgs(1),
	RunE:  runVpcPrefixDelete,
}

func init() {
	vpcPrefixListCmd.Flags().Bool("json", false, "output raw JSON")
	vpcPrefixListCmd.Flags().String("site-id", "", "filter by site ID")
	vpcPrefixListCmd.Flags().String("vpc-id", "", "filter by VPC ID")

	vpcPrefixCreateCmd.Flags().String("name", "", "name for the VPC prefix (required)")
	vpcPrefixCreateCmd.Flags().String("vpc-id", "", "VPC ID (required)")
	vpcPrefixCreateCmd.Flags().Int32("prefix-length", 0, "prefix length (required)")
	vpcPrefixCreateCmd.Flags().String("ip-block-id", "", "IP block ID (optional)")
	vpcPrefixCreateCmd.MarkFlagRequired("name")
	vpcPrefixCreateCmd.MarkFlagRequired("vpc-id")
	vpcPrefixCreateCmd.MarkFlagRequired("prefix-length")

	vpcPrefixUpdateCmd.Flags().String("name", "", "new name (required)")
	vpcPrefixUpdateCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(vpcPrefixCmd)
	vpcPrefixCmd.AddCommand(vpcPrefixListCmd)
	vpcPrefixCmd.AddCommand(vpcPrefixGetCmd)
	vpcPrefixCmd.AddCommand(vpcPrefixCreateCmd)
	vpcPrefixCmd.AddCommand(vpcPrefixUpdateCmd)
	vpcPrefixCmd.AddCommand(vpcPrefixDeleteCmd)
}

func runVpcPrefixList(cmd *cobra.Command, args []string) error {
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

	prefixes, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.VpcPrefix, *http.Response, error) {
		req := apiClient.VPCPrefixAPI.GetAllVpcPrefix(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
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
			return fmt.Errorf("listing VPC prefixes (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing VPC prefixes: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(prefixes))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, prefixes)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, prefixes)
	default:
		return printVpcPrefixTable(os.Stdout, prefixes)
	}
}

func runVpcPrefixGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	prefix, resp, err := apiClient.VPCPrefixAPI.GetVpcPrefix(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting VPC prefix (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting VPC prefix: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, prefix)
	default:
		return printJSON(os.Stdout, prefix)
	}
}

func runVpcPrefixCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	name, _ := cmd.Flags().GetString("name")
	vpcID, _ := cmd.Flags().GetString("vpc-id")
	prefixLen, _ := cmd.Flags().GetInt32("prefix-length")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewVpcPrefixCreateRequest(name, vpcID, prefixLen)
	if cmd.Flags().Changed("ip-block-id") {
		ipBlockID, _ := cmd.Flags().GetString("ip-block-id")
		req.SetIpBlockId(ipBlockID)
	}

	prefix, resp, err := apiClient.VPCPrefixAPI.CreateVpcPrefix(ctx, org).VpcPrefixCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating VPC prefix (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating VPC prefix: %v", err)
	}

	fmt.Fprintf(os.Stderr, "VPC prefix created: %s (%s)\n", ptrStr(prefix.Name), ptrStr(prefix.Id))
	return printJSON(os.Stdout, prefix)
}

func runVpcPrefixUpdate(cmd *cobra.Command, args []string) error {
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

	req := client.NewVpcPrefixUpdateRequest(name)

	prefix, resp, err := apiClient.VPCPrefixAPI.UpdateVpcPrefix(ctx, org, args[0]).VpcPrefixUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating VPC prefix (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating VPC prefix: %v", err)
	}

	fmt.Fprintf(os.Stderr, "VPC prefix updated: %s (%s)\n", ptrStr(prefix.Name), ptrStr(prefix.Id))
	return printJSON(os.Stdout, prefix)
}

func runVpcPrefixDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.VPCPrefixAPI.DeleteVpcPrefix(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting VPC prefix (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting VPC prefix: %v", err)
	}

	fmt.Fprintf(os.Stderr, "VPC prefix deleted: %s\n", args[0])
	return nil
}

func printVpcPrefixTable(w io.Writer, prefixes []client.VpcPrefix) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tVPC ID\tAGE\tID")

	for _, p := range prefixes {
		name := ptrStr(p.Name)
		id := ptrStr(p.Id)
		vpcID := ptrStr(p.VpcId)
		status := ""
		if p.Status != nil {
			status = string(*p.Status)
		}
		age := "<unknown>"
		if p.Created != nil {
			age = formatAge(time.Since(*p.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", name, status, vpcID, age, id)
	}

	return tw.Flush()
}
