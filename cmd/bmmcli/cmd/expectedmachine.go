// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/nvidia/bare-metal-manager-rest/client"
	"github.com/nvidia/bare-metal-manager-rest/cmd/bmmcli/internal/pagination"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var expectedMachineCmd = &cobra.Command{
	Use:   "expected-machine",
	Short: "Expected machine operations",
}

var expectedMachineListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all expected machines",
	RunE:  runExpectedMachineList,
}

var expectedMachineGetCmd = &cobra.Command{
	Use:   "get <expected-machine-id>",
	Short: "Get expected machine details",
	Args:  cobra.ExactArgs(1),
	RunE:  runExpectedMachineGet,
}

var expectedMachineCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an expected machine",
	RunE:  runExpectedMachineCreate,
}

var expectedMachineDeleteCmd = &cobra.Command{
	Use:   "delete <expected-machine-id>",
	Short: "Delete an expected machine",
	Args:  cobra.ExactArgs(1),
	RunE:  runExpectedMachineDelete,
}

func init() {
	expectedMachineListCmd.Flags().Bool("json", false, "output raw JSON")

	expectedMachineCreateCmd.Flags().String("site-id", "", "site ID (required)")
	expectedMachineCreateCmd.Flags().String("bmc-mac-address", "", "BMC MAC address (required)")
	expectedMachineCreateCmd.Flags().String("chassis-serial-number", "", "chassis serial number (required)")
	expectedMachineCreateCmd.MarkFlagRequired("site-id")
	expectedMachineCreateCmd.MarkFlagRequired("bmc-mac-address")
	expectedMachineCreateCmd.MarkFlagRequired("chassis-serial-number")

	rootCmd.AddCommand(expectedMachineCmd)
	expectedMachineCmd.AddCommand(expectedMachineListCmd)
	expectedMachineCmd.AddCommand(expectedMachineGetCmd)
	expectedMachineCmd.AddCommand(expectedMachineCreateCmd)
	expectedMachineCmd.AddCommand(expectedMachineDeleteCmd)
}

func runExpectedMachineList(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	machines, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.ExpectedMachine, *http.Response, error) {
		return apiClient.ExpectedMachineAPI.GetAllExpectedMachine(ctx, org).PageNumber(pageNumber).PageSize(pageSize).Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing expected machines (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing expected machines: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(machines))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, machines)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, machines)
	default:
		return printExpectedMachineTable(os.Stdout, machines)
	}
}

func runExpectedMachineGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	machine, resp, err := apiClient.ExpectedMachineAPI.GetExpectedMachine(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting expected machine (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting expected machine: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, machine)
	default:
		return printJSON(os.Stdout, machine)
	}
}

func runExpectedMachineCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	siteID, _ := cmd.Flags().GetString("site-id")
	bmcMAC, _ := cmd.Flags().GetString("bmc-mac-address")
	chassisSN, _ := cmd.Flags().GetString("chassis-serial-number")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewExpectedMachineCreateRequest(siteID, bmcMAC, chassisSN)

	machine, resp, err := apiClient.ExpectedMachineAPI.CreateExpectedMachine(ctx, org).ExpectedMachineCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating expected machine (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating expected machine: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Expected machine created: %s\n", ptrStr(machine.Id))
	return printJSON(os.Stdout, machine)
}

func runExpectedMachineDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.ExpectedMachineAPI.DeleteExpectedMachine(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting expected machine (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting expected machine: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Expected machine deleted: %s\n", args[0])
	return nil
}

func printExpectedMachineTable(w io.Writer, machines []client.ExpectedMachine) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "SITE ID\tBMC MAC\tCHASSIS SN\tID")

	for _, m := range machines {
		id := ptrStr(m.Id)
		siteID := ptrStr(m.SiteId)
		bmcMAC := ptrStr(m.BmcMacAddress)
		chassisSN := ptrStr(m.ChassisSerialNumber)

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", siteID, bmcMAC, chassisSN, id)
	}

	return tw.Flush()
}
