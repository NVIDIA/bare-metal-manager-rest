// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// admin.go contains read-only commands for admin/informational resources:
// audit, metadata, user, service-account, infrastructure-provider, tenant

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

// -- Audit --

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit log operations",
}

var auditListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all audit log entries",
	RunE:  runAuditList,
}

var auditGetCmd = &cobra.Command{
	Use:   "get <audit-entry-id>",
	Short: "Get audit log entry details",
	Args:  cobra.ExactArgs(1),
	RunE:  runAuditGet,
}

// -- Metadata --

var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "API metadata",
}

var metadataGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get API metadata (version, build time)",
	RunE:  runMetadataGet,
}

// -- User --

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User operations",
}

var userGetCmd = &cobra.Command{
	Use:   "current",
	Short: "Get current user",
	RunE:  runUserGet,
}

// -- Service Account --

var serviceAccountCmd = &cobra.Command{
	Use:   "service-account",
	Short: "Service account operations",
}

var serviceAccountGetCmd = &cobra.Command{
	Use:   "current",
	Short: "Get current service account status",
	RunE:  runServiceAccountGet,
}

// -- Infrastructure Provider --

var infraProviderCmd = &cobra.Command{
	Use:   "infrastructure-provider",
	Short: "Infrastructure provider operations",
}

var infraProviderGetCmd = &cobra.Command{
	Use:   "current",
	Short: "Get current infrastructure provider",
	RunE:  runInfraProviderGet,
}

var infraProviderStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get current infrastructure provider stats",
	RunE:  runInfraProviderStats,
}

// -- Tenant --

var tenantCmd = &cobra.Command{
	Use:   "tenant",
	Short: "Tenant operations",
}

var tenantGetCmd = &cobra.Command{
	Use:   "current",
	Short: "Get current tenant",
	RunE:  runTenantGet,
}

var tenantStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get current tenant stats",
	RunE:  runTenantStats,
}

// -- Machine Capability --

var machineCapabilityCmd = &cobra.Command{
	Use:   "machine-capability",
	Short: "Machine capability operations",
}

var machineCapabilityListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all machine capabilities",
	RunE:  runMachineCapabilityList,
}

func init() {
	auditListCmd.Flags().Bool("json", false, "output raw JSON")

	rootCmd.AddCommand(auditCmd)
	auditCmd.AddCommand(auditListCmd)
	auditCmd.AddCommand(auditGetCmd)

	rootCmd.AddCommand(metadataCmd)
	metadataCmd.AddCommand(metadataGetCmd)

	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userGetCmd)

	rootCmd.AddCommand(serviceAccountCmd)
	serviceAccountCmd.AddCommand(serviceAccountGetCmd)

	rootCmd.AddCommand(infraProviderCmd)
	infraProviderCmd.AddCommand(infraProviderGetCmd)
	infraProviderCmd.AddCommand(infraProviderStatsCmd)

	rootCmd.AddCommand(tenantCmd)
	tenantCmd.AddCommand(tenantGetCmd)
	tenantCmd.AddCommand(tenantStatsCmd)

	machineCapabilityListCmd.Flags().Bool("json", false, "output raw JSON")
	rootCmd.AddCommand(machineCapabilityCmd)
	machineCapabilityCmd.AddCommand(machineCapabilityListCmd)
}

// -- Audit implementations --

func runAuditList(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	entries, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.AuditEntry, *http.Response, error) {
		return apiClient.AuditAPI.GetAllAuditEntry(ctx, org).PageNumber(pageNumber).PageSize(pageSize).Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing audit entries (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing audit entries: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(entries))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, entries)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, entries)
	default:
		return printAuditTable(os.Stdout, entries)
	}
}

func runAuditGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	entry, resp, err := apiClient.AuditAPI.GetAuditEntry(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting audit entry (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting audit entry: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, entry)
	default:
		return printJSON(os.Stdout, entry)
	}
}

func printAuditTable(w io.Writer, entries []client.AuditEntry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "METHOD\tENDPOINT\tSTATUS CODE\tID")

	for _, e := range entries {
		method := ptrStr(e.Method)
		endpoint := ptrStr(e.Endpoint)
		id := ptrStr(e.Id)
		statusCode := ""
		if e.StatusCode != nil {
			statusCode = fmt.Sprintf("%d", *e.StatusCode)
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", method, endpoint, statusCode, id)
	}

	return tw.Flush()
}

// -- Metadata implementation --

func runMetadataGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	metadata, resp, err := apiClient.MetadataAPI.GetMetadata(ctx, org).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting metadata (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting metadata: %v", err)
	}

	return printJSON(os.Stdout, metadata)
}

// -- User implementation --

func runUserGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	user, resp, err := apiClient.UserAPI.GetUser(ctx, org).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting current user (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting current user: %v", err)
	}

	return printJSON(os.Stdout, user)
}

// -- Service Account implementation --

func runServiceAccountGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	sa, resp, err := apiClient.ServiceAccountAPI.GetCurrentServiceAccount(ctx, org).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting service account (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting service account: %v", err)
	}

	return printJSON(os.Stdout, sa)
}

// -- Infrastructure Provider implementations --

func runInfraProviderGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	provider, resp, err := apiClient.InfrastructureProviderAPI.GetCurrentInfrastructureProvider(ctx, org).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting infrastructure provider (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting infrastructure provider: %v", err)
	}

	return printJSON(os.Stdout, provider)
}

func runInfraProviderStats(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	stats, resp, err := apiClient.InfrastructureProviderAPI.GetCurrentInfrastructureProviderStats(ctx, org).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting infrastructure provider stats (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting infrastructure provider stats: %v", err)
	}

	return printJSON(os.Stdout, stats)
}

// -- Tenant implementations --

func runTenantGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	tenant, resp, err := apiClient.TenantAPI.GetCurrentTenant(ctx, org).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting current tenant (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting current tenant: %v", err)
	}

	return printJSON(os.Stdout, tenant)
}

func runTenantStats(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	stats, resp, err := apiClient.TenantAPI.GetCurrentTenantStats(ctx, org).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting tenant stats (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting tenant stats: %v", err)
	}

	return printJSON(os.Stdout, stats)
}

// -- Machine Capability implementation --

func runMachineCapabilityList(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	caps, resp, err := apiClient.MachineAPI.GetAllMachineCapabilities(ctx, org).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing machine capabilities (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing machine capabilities: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(caps))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, caps)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, caps)
	default:
		return printJSON(os.Stdout, caps)
	}
}
