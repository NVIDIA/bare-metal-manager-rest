// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/nvidia/bare-metal-manager-rest/client"
	"github.com/nvidia/bare-metal-manager-rest/cmd/bmmcli/internal/pagination"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var machineCmd = &cobra.Command{
	Use:   "machine",
	Short: "Machine operations",
}

var machineListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all machines",
	RunE:  runMachineList,
}

var machineGetCmd = &cobra.Command{
	Use:   "get <machine-id>",
	Short: "Get machine details",
	Args:  cobra.ExactArgs(1),
	RunE:  runMachineGet,
}

var machineUpdateCmd = &cobra.Command{
	Use:   "update <machine-id>",
	Short: "Update a machine",
	Args:  cobra.ExactArgs(1),
	RunE:  runMachineUpdate,
}

var machineDeleteCmd = &cobra.Command{
	Use:   "delete <machine-id>",
	Short: "Delete a machine",
	Args:  cobra.ExactArgs(1),
	RunE:  runMachineDelete,
}

func init() {
	machineListCmd.Flags().Bool("json", false, "output raw JSON")
	machineListCmd.Flags().String("site-id", "", "filter by site ID")
	machineListCmd.Flags().String("vpc-id", "", "filter by VPC ID")
	machineListCmd.Flags().String("status", "", "filter by status")
	machineListCmd.Flags().String("instance-type-id", "", "filter by instance type ID")

	machineUpdateCmd.Flags().String("instance-type-id", "", "instance type ID to assign")
	machineUpdateCmd.Flags().Bool("clear-instance-type", false, "clear assigned instance type")
	machineUpdateCmd.Flags().Bool("maintenance-mode", false, "set maintenance mode")
	machineUpdateCmd.Flags().String("maintenance-message", "", "maintenance message")

	rootCmd.AddCommand(machineCmd)
	machineCmd.AddCommand(machineListCmd)
	machineCmd.AddCommand(machineGetCmd)
	machineCmd.AddCommand(machineUpdateCmd)
	machineCmd.AddCommand(machineDeleteCmd)
}

func runMachineList(cmd *cobra.Command, args []string) error {
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
	status, _ := cmd.Flags().GetString("status")
	instanceTypeID, _ := cmd.Flags().GetString("instance-type-id")

	machines, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.Machine, *http.Response, error) {
		req := apiClient.MachineAPI.GetAllMachine(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		if status != "" {
			req = req.Status(status)
		}
		if instanceTypeID != "" {
			req = req.InstanceTypeId(instanceTypeID)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing machines (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing machines: %v", err)
	}

	machineVpcNames, machineIDsInScope, vpcErr := fetchMachineVpcNames(apiClient, ctx, org, siteID, vpcID)
	if vpcErr != nil {
		if vpcID != "" {
			return vpcErr
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Note: could not load machine VPC names: %v\n", vpcErr)
	}

	summaryResp := resp
	if vpcID != "" {
		filtered := make([]client.Machine, 0, len(machines))
		for _, m := range machines {
			id := strings.TrimSpace(ptrStr(m.Id))
			if _, ok := machineIDsInScope[id]; ok {
				filtered = append(filtered, m)
			}
		}
		machines = filtered
		summaryResp = nil
	}

	pagination.PrintSummary(cmd.ErrOrStderr(), summaryResp, len(machines))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, machines)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, machines)
	default:
		return printMachineTable(os.Stdout, machines, machineVpcNames)
	}
}

func fetchMachineVpcNames(apiClient *client.APIClient, ctx context.Context, org, siteID, vpcID string) (map[string]string, map[string]struct{}, error) {
	instances, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.Instance, *http.Response, error) {
		req := apiClient.InstanceAPI.GetAllInstance(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
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
			return nil, nil, fmt.Errorf("listing instances for machine VPC names (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return nil, nil, fmt.Errorf("listing instances for machine VPC names: %v", err)
	}

	machineIDSet := make(map[string]struct{})
	vpcIDsByMachineID := make(map[string]map[string]struct{})
	for _, inst := range instances {
		machineID := ""
		if machineIDPtr, ok := inst.GetMachineIdOk(); ok && machineIDPtr != nil {
			machineID = strings.TrimSpace(*machineIDPtr)
		}
		instanceVpcID := strings.TrimSpace(ptrStr(inst.VpcId))
		if machineID == "" || instanceVpcID == "" {
			continue
		}
		machineIDSet[machineID] = struct{}{}
		if vpcIDsByMachineID[machineID] == nil {
			vpcIDsByMachineID[machineID] = make(map[string]struct{})
		}
		vpcIDsByMachineID[machineID][instanceVpcID] = struct{}{}
	}

	vpcNamesByID := make(map[string]string)
	vpcs, _, vpcErr := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.VPC, *http.Response, error) {
		req := apiClient.VPCAPI.GetAllVpc(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
	})
	if vpcErr == nil {
		for _, v := range vpcs {
			id := strings.TrimSpace(ptrStr(v.Id))
			if id == "" {
				continue
			}
			name := strings.TrimSpace(ptrStr(v.Name))
			if name == "" {
				name = id
			}
			vpcNamesByID[id] = name
		}
	}

	vpcNamesByMachineID := make(map[string]string, len(vpcIDsByMachineID))
	for machineID, vpcSet := range vpcIDsByMachineID {
		names := make([]string, 0, len(vpcSet))
		for id := range vpcSet {
			name := strings.TrimSpace(vpcNamesByID[id])
			if name == "" {
				name = id
			}
			names = append(names, name)
		}
		sort.Strings(names)
		vpcNamesByMachineID[machineID] = strings.Join(names, ",")
	}

	return vpcNamesByMachineID, machineIDSet, nil
}

func runMachineGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	machine, resp, err := apiClient.MachineAPI.GetMachine(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting machine (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting machine: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, machine)
	default:
		return printJSON(os.Stdout, machine)
	}
}

func runMachineUpdate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewMachineUpdateRequest()
	if cmd.Flags().Changed("instance-type-id") {
		itID, _ := cmd.Flags().GetString("instance-type-id")
		req.SetInstanceTypeId(itID)
	}
	if cmd.Flags().Changed("clear-instance-type") {
		v, _ := cmd.Flags().GetBool("clear-instance-type")
		req.SetClearInstanceType(v)
	}
	if cmd.Flags().Changed("maintenance-mode") {
		v, _ := cmd.Flags().GetBool("maintenance-mode")
		req.SetSetMaintenanceMode(v)
	}
	if cmd.Flags().Changed("maintenance-message") {
		msg, _ := cmd.Flags().GetString("maintenance-message")
		req.SetMaintenanceMessage(msg)
	}

	machine, resp, err := apiClient.MachineAPI.UpdateMachine(ctx, org, args[0]).MachineUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating machine (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating machine: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Machine updated: %s\n", ptrStr(machine.Id))
	return printJSON(os.Stdout, machine)
}

func runMachineDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.MachineAPI.DeleteMachine(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting machine (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting machine: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Machine deleted: %s\n", args[0])
	return nil
}

func printMachineTable(w io.Writer, machines []client.Machine, vpcNamesByMachineID map[string]string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "STATUS\tSITE ID\tVPC\tAGE\tID")

	for _, m := range machines {
		id := ptrStr(m.Id)
		siteID := ptrStr(m.SiteId)
		vpcNames := strings.TrimSpace(vpcNamesByMachineID[id])
		if vpcNames == "" {
			vpcNames = "-"
		}
		status := ""
		if m.Status != nil {
			status = string(*m.Status)
		}
		age := "<unknown>"
		if m.Created != nil {
			age = formatAge(time.Since(*m.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", status, siteID, vpcNames, age, id)
	}

	return tw.Flush()
}
