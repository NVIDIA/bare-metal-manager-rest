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

var dpuExtSvcCmd = &cobra.Command{
	Use:   "dpu-extension-service",
	Short: "DPU extension service operations",
}

var dpuExtSvcListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all DPU extension services",
	RunE:  runDPUExtSvcList,
}

var dpuExtSvcGetCmd = &cobra.Command{
	Use:   "get <service-id>",
	Short: "Get DPU extension service details",
	Args:  cobra.ExactArgs(1),
	RunE:  runDPUExtSvcGet,
}

var dpuExtSvcCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a DPU extension service",
	RunE:  runDPUExtSvcCreate,
}

var dpuExtSvcUpdateCmd = &cobra.Command{
	Use:   "update <service-id>",
	Short: "Update a DPU extension service",
	Args:  cobra.ExactArgs(1),
	RunE:  runDPUExtSvcUpdate,
}

var dpuExtSvcDeleteCmd = &cobra.Command{
	Use:   "delete <service-id>",
	Short: "Delete a DPU extension service",
	Args:  cobra.ExactArgs(1),
	RunE:  runDPUExtSvcDelete,
}

func init() {
	dpuExtSvcListCmd.Flags().Bool("json", false, "output raw JSON")
	dpuExtSvcListCmd.Flags().String("site-id", "", "filter by site ID")

	dpuExtSvcCreateCmd.Flags().String("name", "", "name for the service (required)")
	dpuExtSvcCreateCmd.Flags().String("service-type", "", "service type (required)")
	dpuExtSvcCreateCmd.Flags().String("site-id", "", "site ID (required)")
	dpuExtSvcCreateCmd.Flags().String("data", "", "service data (required)")
	dpuExtSvcCreateCmd.MarkFlagRequired("name")
	dpuExtSvcCreateCmd.MarkFlagRequired("service-type")
	dpuExtSvcCreateCmd.MarkFlagRequired("site-id")
	dpuExtSvcCreateCmd.MarkFlagRequired("data")

	dpuExtSvcUpdateCmd.Flags().String("name", "", "new name")
	dpuExtSvcUpdateCmd.Flags().String("description", "", "new description")

	rootCmd.AddCommand(dpuExtSvcCmd)
	dpuExtSvcCmd.AddCommand(dpuExtSvcListCmd)
	dpuExtSvcCmd.AddCommand(dpuExtSvcGetCmd)
	dpuExtSvcCmd.AddCommand(dpuExtSvcCreateCmd)
	dpuExtSvcCmd.AddCommand(dpuExtSvcUpdateCmd)
	dpuExtSvcCmd.AddCommand(dpuExtSvcDeleteCmd)
}

func runDPUExtSvcList(cmd *cobra.Command, args []string) error {
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

	services, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.DpuExtensionService, *http.Response, error) {
		req := apiClient.DPUExtensionServiceAPI.GetAllDpuExtensionService(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing DPU extension services (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing DPU extension services: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(services))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, services)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, services)
	default:
		return printDPUExtSvcTable(os.Stdout, services)
	}
}

func runDPUExtSvcGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	svc, resp, err := apiClient.DPUExtensionServiceAPI.GetDpuExtensionService(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting DPU extension service (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting DPU extension service: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, svc)
	default:
		return printJSON(os.Stdout, svc)
	}
}

func runDPUExtSvcCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	name, _ := cmd.Flags().GetString("name")
	serviceType, _ := cmd.Flags().GetString("service-type")
	siteID, _ := cmd.Flags().GetString("site-id")
	data, _ := cmd.Flags().GetString("data")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewDpuExtensionServiceCreateRequest(name, serviceType, siteID, data)

	svc, resp, err := apiClient.DPUExtensionServiceAPI.CreateDpuExtensionService(ctx, org).DpuExtensionServiceCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating DPU extension service (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating DPU extension service: %v", err)
	}

	fmt.Fprintf(os.Stderr, "DPU extension service created: %s (%s)\n", ptrStr(svc.Name), ptrStr(svc.Id))
	return printJSON(os.Stdout, svc)
}

func runDPUExtSvcUpdate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewDpuExtensionServiceUpdateRequest()
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		req.SetName(name)
	}
	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		req.SetDescription(description)
	}

	svc, resp, err := apiClient.DPUExtensionServiceAPI.UpdateDpuExtensionService(ctx, org, args[0]).DpuExtensionServiceUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating DPU extension service (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating DPU extension service: %v", err)
	}

	fmt.Fprintf(os.Stderr, "DPU extension service updated: %s (%s)\n", ptrStr(svc.Name), ptrStr(svc.Id))
	return printJSON(os.Stdout, svc)
}

func runDPUExtSvcDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.DPUExtensionServiceAPI.DeleteDpuExtensionService(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting DPU extension service (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting DPU extension service: %v", err)
	}

	fmt.Fprintf(os.Stderr, "DPU extension service deleted: %s\n", args[0])
	return nil
}

func printDPUExtSvcTable(w io.Writer, services []client.DpuExtensionService) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tTYPE\tSITE ID\tAGE\tID")

	for _, s := range services {
		name := ptrStr(s.Name)
		id := ptrStr(s.Id)
		siteID := ptrStr(s.SiteId)
		svcType := ptrStr(s.ServiceType)
		age := "<unknown>"
		if s.Created != nil {
			age = formatAge(time.Since(*s.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", name, svcType, siteID, age, id)
	}

	return tw.Flush()
}
