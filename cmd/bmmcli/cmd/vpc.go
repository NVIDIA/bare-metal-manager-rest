// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/nvidia/bare-metal-manager-rest/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var vpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "VPC operations",
}

var vpcListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VPCs",
	RunE:  runVpcList,
}

var vpcCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a VPC",
	RunE:  runVpcCreate,
}

func init() {
	vpcListCmd.Flags().Bool("json", false, "output raw JSON")

	vpcCreateCmd.Flags().String("name", "", "name for the VPC (required)")
	vpcCreateCmd.Flags().String("site-id", "", "site ID where the VPC should be created (required)")
	vpcCreateCmd.Flags().String("description", "", "optional description")
	vpcCreateCmd.MarkFlagRequired("name")
	vpcCreateCmd.MarkFlagRequired("site-id")

	rootCmd.AddCommand(vpcCmd)
	vpcCmd.AddCommand(vpcListCmd)
	vpcCmd.AddCommand(vpcCreateCmd)
}

func runVpcList(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	vpcs, resp, err := apiClient.VPCAPI.GetAllVpc(ctx, org).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing VPCs (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing VPCs: %v", err)
	}

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, vpcs)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, vpcs)
	default:
		return printVpcTable(os.Stdout, vpcs)
	}
}

func runVpcCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	name, _ := cmd.Flags().GetString("name")
	siteID, _ := cmd.Flags().GetString("site-id")
	description, _ := cmd.Flags().GetString("description")

	req := client.NewVpcCreateRequest(name, siteID)
	if description != "" {
		req.SetDescription(description)
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	vpc, resp, err := apiClient.VPCAPI.CreateVpc(ctx, org).VpcCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating VPC (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating VPC: %v", err)
	}

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, vpc)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, vpc)
	default:
		fmt.Fprintf(os.Stderr, "VPC created: %s (%s)\n", ptrStr(vpc.Name), ptrStr(vpc.Id))
		return printVpcTable(os.Stdout, []client.VPC{*vpc})
	}
}

func printVpcTable(w io.Writer, vpcs []client.VPC) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE ID\tAGE\tID")

	for _, v := range vpcs {
		name := ptrStr(v.Name)
		id := ptrStr(v.Id)
		siteID := ptrStr(v.SiteId)
		status := ""
		if v.Status != nil {
			status = string(*v.Status)
		}
		age := "<unknown>"
		if v.Created != nil {
			age = formatAge(time.Since(*v.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", name, status, siteID, age, id)
	}

	return tw.Flush()
}

// tryReadBody reads and returns the response body as a string for error reporting
func tryReadBody(body io.ReadCloser) string {
	if body == nil {
		return ""
	}
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return ""
	}
	return string(data)
}
