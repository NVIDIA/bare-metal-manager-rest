// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/nvidia/bare-metal-manager-rest/client"
	"github.com/nvidia/bare-metal-manager-rest/cmd/bmmcli/internal/pagination"
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

var vpcGetCmd = &cobra.Command{
	Use:   "get <vpc-id>",
	Short: "Get VPC details",
	Args:  cobra.ExactArgs(1),
	RunE:  runVpcGet,
}

var vpcUpdateCmd = &cobra.Command{
	Use:   "update <vpc-id>",
	Short: "Update a VPC",
	Args:  cobra.ExactArgs(1),
	RunE:  runVpcUpdate,
}

var vpcDeleteCmd = &cobra.Command{
	Use:   "delete <vpc-id>",
	Short: "Delete a VPC",
	Args:  cobra.ExactArgs(1),
	RunE:  runVpcDelete,
}

func init() {
	vpcListCmd.Flags().Bool("json", false, "output raw JSON")
	vpcListCmd.Flags().String("site-id", "", "filter by site ID")

	vpcCreateCmd.Flags().String("name", "", "name for the VPC (required)")
	vpcCreateCmd.Flags().String("site-id", "", "site ID where the VPC should be created (required)")
	vpcCreateCmd.Flags().String("description", "", "optional description")
	vpcCreateCmd.MarkFlagRequired("name")
	vpcCreateCmd.MarkFlagRequired("site-id")

	vpcUpdateCmd.Flags().String("name", "", "new name")
	vpcUpdateCmd.Flags().String("description", "", "new description")

	rootCmd.AddCommand(vpcCmd)
	vpcCmd.AddCommand(vpcListCmd)
	vpcCmd.AddCommand(vpcCreateCmd)
	vpcCmd.AddCommand(vpcGetCmd)
	vpcCmd.AddCommand(vpcUpdateCmd)
	vpcCmd.AddCommand(vpcDeleteCmd)
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
	siteID, _ := cmd.Flags().GetString("site-id")

	vpcs, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.VPC, *http.Response, error) {
		req := apiClient.VPCAPI.GetAllVpc(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing VPCs (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing VPCs: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(vpcs))
	usageByVpc, usageErr := fetchVpcUsageCounts(apiClient, ctx, org, siteID)
	if usageErr != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Note: could not load VPC usage counts: %v\n", usageErr)
	}
	siteNamesByID, siteErr := fetchSiteNamesByID(apiClient, ctx, org)
	if siteErr != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Note: could not load site names: %v\n", siteErr)
	}

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, vpcs)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, vpcs)
	default:
		return printVpcTable(os.Stdout, vpcs, usageByVpc, siteNamesByID)
	}
}

func runVpcGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	vpc, resp, err := apiClient.VPCAPI.GetVpc(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting VPC (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting VPC: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, vpc)
	default:
		return printJSON(os.Stdout, vpc)
	}
}

func runVpcUpdate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewVpcUpdateRequest()
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		req.SetName(name)
	}
	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		req.SetDescription(description)
	}

	vpc, resp, err := apiClient.VPCAPI.UpdateVpc(ctx, org, args[0]).VpcUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating VPC (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating VPC: %v", err)
	}

	fmt.Fprintf(os.Stderr, "VPC updated: %s (%s)\n", ptrStr(vpc.Name), ptrStr(vpc.Id))
	return printJSON(os.Stdout, vpc)
}

func runVpcDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.VPCAPI.DeleteVpc(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting VPC (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting VPC: %v", err)
	}

	fmt.Fprintf(os.Stderr, "VPC deleted: %s\n", args[0])
	return nil
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
		return printVpcTable(os.Stdout, []client.VPC{*vpc}, nil, nil)
	}
}

type vpcUsageCounts struct {
	Instances int
	Machines  int
}

func fetchVpcUsageCounts(apiClient *client.APIClient, ctx context.Context, org, siteID string) (map[string]vpcUsageCounts, error) {
	instances, _, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.Instance, *http.Response, error) {
		req := apiClient.InstanceAPI.GetAllInstance(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
	})
	if err != nil {
		return nil, err
	}

	usageByVpc := make(map[string]vpcUsageCounts)
	machineSetByVpc := make(map[string]map[string]struct{})
	for _, inst := range instances {
		vpcID := strings.TrimSpace(ptrStr(inst.VpcId))
		if vpcID == "" {
			continue
		}

		usage := usageByVpc[vpcID]
		usage.Instances++
		usageByVpc[vpcID] = usage

		if machineID, ok := inst.GetMachineIdOk(); ok && machineID != nil {
			id := strings.TrimSpace(*machineID)
			if id != "" {
				if machineSetByVpc[vpcID] == nil {
					machineSetByVpc[vpcID] = make(map[string]struct{})
				}
				machineSetByVpc[vpcID][id] = struct{}{}
			}
		}
	}

	for vpcID, set := range machineSetByVpc {
		usage := usageByVpc[vpcID]
		usage.Machines = len(set)
		usageByVpc[vpcID] = usage
	}
	return usageByVpc, nil
}

func fetchSiteNamesByID(apiClient *client.APIClient, ctx context.Context, org string) (map[string]string, error) {
	sites, _, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.Site, *http.Response, error) {
		return apiClient.SiteAPI.GetAllSite(ctx, org).PageNumber(pageNumber).PageSize(pageSize).Execute()
	})
	if err != nil {
		return nil, err
	}

	siteNamesByID := make(map[string]string, len(sites))
	for _, s := range sites {
		id := strings.TrimSpace(ptrStr(s.Id))
		if id == "" {
			continue
		}
		siteNamesByID[id] = strings.TrimSpace(ptrStr(s.Name))
	}
	return siteNamesByID, nil
}

func printVpcTable(w io.Writer, vpcs []client.VPC, usageByVpc map[string]vpcUsageCounts, siteNamesByID map[string]string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE\tINSTANCES\tMACHINES\tAGE\tID")

	for _, v := range vpcs {
		name := ptrStr(v.Name)
		id := ptrStr(v.Id)
		siteID := ptrStr(v.SiteId)
		siteName := strings.TrimSpace(siteNamesByID[siteID])
		if siteName == "" {
			siteName = siteID
		}
		status := ""
		if v.Status != nil {
			status = string(*v.Status)
		}
		age := "<unknown>"
		if v.Created != nil {
			age = formatAge(time.Since(*v.Created))
		}

		usage := usageByVpc[id]
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%d\t%s\t%s\n", name, status, siteName, usage.Instances, usage.Machines, age, id)
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
