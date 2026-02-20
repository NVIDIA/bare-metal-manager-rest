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

var nsgCmd = &cobra.Command{
	Use:   "network-security-group",
	Short: "Network security group operations",
}

var nsgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all network security groups",
	RunE:  runNSGListCmd,
}

var nsgGetCmd = &cobra.Command{
	Use:   "get <nsg-id>",
	Short: "Get network security group details",
	Args:  cobra.ExactArgs(1),
	RunE:  runNSGGetCmd,
}

var nsgDeleteCmd = &cobra.Command{
	Use:   "delete <nsg-id>",
	Short: "Delete a network security group",
	Args:  cobra.ExactArgs(1),
	RunE:  runNSGDeleteCmd,
}

func init() {
	nsgListCmd.Flags().Bool("json", false, "output raw JSON")
	nsgListCmd.Flags().String("site-id", "", "filter by site ID")

	rootCmd.AddCommand(nsgCmd)
	nsgCmd.AddCommand(nsgListCmd)
	nsgCmd.AddCommand(nsgGetCmd)
	nsgCmd.AddCommand(nsgDeleteCmd)
}

func runNSGListCmd(cmd *cobra.Command, args []string) error {
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

	nsgs, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.NetworkSecurityGroup, *http.Response, error) {
		req := apiClient.NetworkSecurityGroupAPI.GetAllNetworkSecurityGroup(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing network security groups (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing network security groups: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(nsgs))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, nsgs)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, nsgs)
	default:
		return printNSGTable(os.Stdout, nsgs)
	}
}

func runNSGGetCmd(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	nsg, resp, err := apiClient.NetworkSecurityGroupAPI.GetNetworkSecurityGroup(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting network security group (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting network security group: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, nsg)
	default:
		return printJSON(os.Stdout, nsg)
	}
}

func runNSGDeleteCmd(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.NetworkSecurityGroupAPI.DeleteNetworkSecurityGroup(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting network security group (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting network security group: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Network security group deleted: %s\n", args[0])
	return nil
}

func printNSGTable(w io.Writer, nsgs []client.NetworkSecurityGroup) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tAGE\tID")

	for _, n := range nsgs {
		name := ptrStr(n.Name)
		id := ptrStr(n.Id)
		status := ""
		if n.Status != nil {
			status = string(*n.Status)
		}
		age := "<unknown>"
		if n.Created != nil {
			age = formatAge(time.Since(*n.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", name, status, age, id)
	}

	return tw.Flush()
}
