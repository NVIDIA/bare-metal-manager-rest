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

var operatingSystemCmd = &cobra.Command{
	Use:   "operating-system",
	Short: "Operating system operations",
}

var operatingSystemListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all operating systems",
	RunE:  runOperatingSystemList,
}

var operatingSystemGetCmd = &cobra.Command{
	Use:   "get <os-id>",
	Short: "Get operating system details",
	Args:  cobra.ExactArgs(1),
	RunE:  runOperatingSystemGet,
}

var operatingSystemDeleteCmd = &cobra.Command{
	Use:   "delete <os-id>",
	Short: "Delete an operating system",
	Args:  cobra.ExactArgs(1),
	RunE:  runOperatingSystemDelete,
}

func init() {
	operatingSystemListCmd.Flags().Bool("json", false, "output raw JSON")

	rootCmd.AddCommand(operatingSystemCmd)
	operatingSystemCmd.AddCommand(operatingSystemListCmd)
	operatingSystemCmd.AddCommand(operatingSystemGetCmd)
	operatingSystemCmd.AddCommand(operatingSystemDeleteCmd)
}

func runOperatingSystemList(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	osList, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.OperatingSystem, *http.Response, error) {
		return apiClient.OperatingSystemAPI.GetAllOperatingSystem(ctx, org).PageNumber(pageNumber).PageSize(pageSize).Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing operating systems (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing operating systems: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(osList))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, osList)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, osList)
	default:
		return printOperatingSystemTable(os.Stdout, osList)
	}
}

func runOperatingSystemGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	osItem, resp, err := apiClient.OperatingSystemAPI.GetOperatingSystem(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting operating system (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting operating system: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, osItem)
	default:
		return printJSON(os.Stdout, osItem)
	}
}

func runOperatingSystemDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.OperatingSystemAPI.DeleteOperatingSystem(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting operating system (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting operating system: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Operating system deleted: %s\n", args[0])
	return nil
}

func printOperatingSystemTable(w io.Writer, osList []client.OperatingSystem) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tAGE\tID")

	for _, o := range osList {
		name := ptrStr(o.Name)
		id := ptrStr(o.Id)
		status := ""
		if o.Status != nil {
			status = string(*o.Status)
		}
		age := "<unknown>"
		if o.Created != nil {
			age = formatAge(time.Since(*o.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", name, status, age, id)
	}

	return tw.Flush()
}
