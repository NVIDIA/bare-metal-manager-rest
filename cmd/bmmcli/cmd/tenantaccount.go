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

var tenantAccountCmd = &cobra.Command{
	Use:   "tenant-account",
	Short: "Tenant account operations",
}

var tenantAccountListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tenant accounts",
	RunE:  runTenantAccountList,
}

var tenantAccountGetCmd = &cobra.Command{
	Use:   "get <account-id>",
	Short: "Get tenant account details",
	Args:  cobra.ExactArgs(1),
	RunE:  runTenantAccountGet,
}

var tenantAccountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a tenant account",
	RunE:  runTenantAccountCreate,
}

var tenantAccountDeleteCmd = &cobra.Command{
	Use:   "delete <account-id>",
	Short: "Delete a tenant account",
	Args:  cobra.ExactArgs(1),
	RunE:  runTenantAccountDelete,
}

func init() {
	tenantAccountListCmd.Flags().Bool("json", false, "output raw JSON")

	tenantAccountCreateCmd.Flags().String("infrastructure-provider-id", "", "infrastructure provider ID (required)")
	tenantAccountCreateCmd.Flags().String("tenant-org", "", "tenant organization (required)")
	tenantAccountCreateCmd.MarkFlagRequired("infrastructure-provider-id")
	tenantAccountCreateCmd.MarkFlagRequired("tenant-org")

	rootCmd.AddCommand(tenantAccountCmd)
	tenantAccountCmd.AddCommand(tenantAccountListCmd)
	tenantAccountCmd.AddCommand(tenantAccountGetCmd)
	tenantAccountCmd.AddCommand(tenantAccountCreateCmd)
	tenantAccountCmd.AddCommand(tenantAccountDeleteCmd)
}

func runTenantAccountList(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	accounts, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.TenantAccount, *http.Response, error) {
		return apiClient.TenantAccountAPI.GetAllTenantAccount(ctx, org).PageNumber(pageNumber).PageSize(pageSize).Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing tenant accounts (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing tenant accounts: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(accounts))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, accounts)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, accounts)
	default:
		return printTenantAccountTable(os.Stdout, accounts)
	}
}

func runTenantAccountGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	account, resp, err := apiClient.TenantAccountAPI.GetTenantAccount(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting tenant account (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting tenant account: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, account)
	default:
		return printJSON(os.Stdout, account)
	}
}

func runTenantAccountCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	ipID, _ := cmd.Flags().GetString("infrastructure-provider-id")
	tenantOrg, _ := cmd.Flags().GetString("tenant-org")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewTenantAccountCreateRequest(ipID, tenantOrg)

	account, resp, err := apiClient.TenantAccountAPI.CreateTenantAccount(ctx, org).TenantAccountCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating tenant account (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating tenant account: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Tenant account created: %s\n", ptrStr(account.Id))
	return printJSON(os.Stdout, account)
}

func runTenantAccountDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.TenantAccountAPI.DeleteTenantAccount(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting tenant account (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting tenant account: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Tenant account deleted: %s\n", args[0])
	return nil
}

func printTenantAccountTable(w io.Writer, accounts []client.TenantAccount) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "STATUS\tTENANT ORG\tINFRA PROVIDER ID\tAGE\tID")

	for _, a := range accounts {
		id := ptrStr(a.Id)
		ipID := ptrStr(a.InfrastructureProviderId)
		tenantOrg := ""
		if a.TenantOrg.IsSet() {
			tenantOrg = *a.TenantOrg.Get()
		}
		status := ""
		if a.Status != nil {
			status = string(*a.Status)
		}
		age := "<unknown>"
		if a.Created != nil {
			age = formatAge(time.Since(*a.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", status, tenantOrg, ipID, age, id)
	}

	return tw.Flush()
}
