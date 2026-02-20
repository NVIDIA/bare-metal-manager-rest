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

var operatingSystemCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an operating system",
	RunE:  runOperatingSystemCreate,
}

var operatingSystemUpdateCmd = &cobra.Command{
	Use:   "update <os-id>",
	Short: "Update an operating system",
	Args:  cobra.ExactArgs(1),
	RunE:  runOperatingSystemUpdate,
}

var operatingSystemDeleteCmd = &cobra.Command{
	Use:   "delete <os-id>",
	Short: "Delete an operating system",
	Args:  cobra.ExactArgs(1),
	RunE:  runOperatingSystemDelete,
}

func init() {
	operatingSystemListCmd.Flags().Bool("json", false, "output raw JSON")
	operatingSystemListCmd.Flags().String("site-id", "", "filter by site ID")

	operatingSystemCreateCmd.Flags().String("name", "", "name for the operating system (required)")
	operatingSystemCreateCmd.Flags().String("description", "", "optional description")
	operatingSystemCreateCmd.Flags().String("tenant-id", "", "tenant ID (required unless --infrastructure-provider-id is provided)")
	operatingSystemCreateCmd.Flags().String("infrastructure-provider-id", "", "infrastructure provider ID (required unless --tenant-id is provided)")
	operatingSystemCreateCmd.Flags().StringSlice("site-id", nil, "site ID(s) (can be repeated)")
	operatingSystemCreateCmd.Flags().String("ipxe-script", "", "iPXE script or URL (for iPXE-based OS)")
	operatingSystemCreateCmd.Flags().String("image-url", "", "image URL (for image-based OS)")
	operatingSystemCreateCmd.Flags().String("image-sha", "", "image SHA hash")
	operatingSystemCreateCmd.Flags().String("image-auth-type", "", "image auth type (basic, bearer, token)")
	operatingSystemCreateCmd.Flags().String("image-auth-token", "", "image auth token")
	operatingSystemCreateCmd.Flags().String("image-disk", "", "target disk path for image OS")
	operatingSystemCreateCmd.Flags().String("root-fs-id", "", "root filesystem UUID")
	operatingSystemCreateCmd.Flags().String("root-fs-label", "", "root filesystem label")
	operatingSystemCreateCmd.Flags().String("user-data", "", "user-data payload")
	operatingSystemCreateCmd.Flags().Bool("is-cloud-init", false, "mark as cloud-init based OS")
	operatingSystemCreateCmd.Flags().Bool("allow-override", false, "allow user-data override at instance creation")
	operatingSystemCreateCmd.Flags().Bool("phone-home-enabled", false, "enable phone-home service")
	operatingSystemCreateCmd.MarkFlagRequired("name")

	operatingSystemUpdateCmd.Flags().String("name", "", "new name")
	operatingSystemUpdateCmd.Flags().String("description", "", "new description")
	operatingSystemUpdateCmd.Flags().String("ipxe-script", "", "new iPXE script or URL")
	operatingSystemUpdateCmd.Flags().String("image-url", "", "new image URL")
	operatingSystemUpdateCmd.Flags().String("image-sha", "", "new image SHA hash")
	operatingSystemUpdateCmd.Flags().String("image-auth-type", "", "new image auth type")
	operatingSystemUpdateCmd.Flags().String("image-auth-token", "", "new image auth token")
	operatingSystemUpdateCmd.Flags().String("image-disk", "", "new target disk path")
	operatingSystemUpdateCmd.Flags().String("root-fs-id", "", "new root filesystem UUID")
	operatingSystemUpdateCmd.Flags().String("root-fs-label", "", "new root filesystem label")
	operatingSystemUpdateCmd.Flags().String("user-data", "", "new user-data payload")
	operatingSystemUpdateCmd.Flags().String("deactivation-note", "", "deactivation note when setting --is-active=false")
	operatingSystemUpdateCmd.Flags().Bool("is-cloud-init", false, "set cloud-init mode")
	operatingSystemUpdateCmd.Flags().Bool("allow-override", false, "set user-data override behavior")
	operatingSystemUpdateCmd.Flags().Bool("phone-home-enabled", false, "set phone-home behavior")
	operatingSystemUpdateCmd.Flags().Bool("is-active", false, "set active/inactive state")

	rootCmd.AddCommand(operatingSystemCmd)
	operatingSystemCmd.AddCommand(operatingSystemListCmd)
	operatingSystemCmd.AddCommand(operatingSystemGetCmd)
	operatingSystemCmd.AddCommand(operatingSystemCreateCmd)
	operatingSystemCmd.AddCommand(operatingSystemUpdateCmd)
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
	siteID, _ := cmd.Flags().GetString("site-id")

	osList, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.OperatingSystem, *http.Response, error) {
		req := apiClient.OperatingSystemAPI.GetAllOperatingSystem(ctx, org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
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

func runOperatingSystemCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	name, _ := cmd.Flags().GetString("name")
	tenantID, _ := cmd.Flags().GetString("tenant-id")
	providerID, _ := cmd.Flags().GetString("infrastructure-provider-id")
	if tenantID == "" && providerID == "" {
		return fmt.Errorf("either --tenant-id or --infrastructure-provider-id is required")
	}
	if tenantID != "" && providerID != "" {
		return fmt.Errorf("only one of --tenant-id or --infrastructure-provider-id can be specified")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewOperatingSystemCreateRequest(name)
	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		req.SetDescription(description)
	}
	if tenantID != "" {
		req.SetTenantId(tenantID)
	}
	if providerID != "" {
		req.SetInfrastructureProviderId(providerID)
	}
	if cmd.Flags().Changed("site-id") {
		siteIDs, _ := cmd.Flags().GetStringSlice("site-id")
		req.SetSiteIds(siteIDs)
	}
	if cmd.Flags().Changed("ipxe-script") {
		ipxeScript, _ := cmd.Flags().GetString("ipxe-script")
		req.SetIpxeScript(ipxeScript)
	}
	if cmd.Flags().Changed("image-url") {
		imageURL, _ := cmd.Flags().GetString("image-url")
		req.SetImageUrl(imageURL)
	}
	if cmd.Flags().Changed("image-sha") {
		imageSHA, _ := cmd.Flags().GetString("image-sha")
		req.SetImageSha(imageSHA)
	}
	if cmd.Flags().Changed("image-auth-type") {
		imageAuthType, _ := cmd.Flags().GetString("image-auth-type")
		req.SetImageAuthType(imageAuthType)
	}
	if cmd.Flags().Changed("image-auth-token") {
		imageAuthToken, _ := cmd.Flags().GetString("image-auth-token")
		req.SetImageAuthToken(imageAuthToken)
	}
	if cmd.Flags().Changed("image-disk") {
		imageDisk, _ := cmd.Flags().GetString("image-disk")
		req.SetImageDisk(imageDisk)
	}
	if cmd.Flags().Changed("root-fs-id") {
		rootFsID, _ := cmd.Flags().GetString("root-fs-id")
		req.SetRootFsId(rootFsID)
	}
	if cmd.Flags().Changed("root-fs-label") {
		rootFsLabel, _ := cmd.Flags().GetString("root-fs-label")
		req.SetRootFsLabel(rootFsLabel)
	}
	if cmd.Flags().Changed("user-data") {
		userData, _ := cmd.Flags().GetString("user-data")
		req.SetUserData(userData)
	}
	if cmd.Flags().Changed("is-cloud-init") {
		isCloudInit, _ := cmd.Flags().GetBool("is-cloud-init")
		req.SetIsCloudInit(isCloudInit)
	}
	if cmd.Flags().Changed("allow-override") {
		allowOverride, _ := cmd.Flags().GetBool("allow-override")
		req.SetAllowOverride(allowOverride)
	}
	if cmd.Flags().Changed("phone-home-enabled") {
		phoneHomeEnabled, _ := cmd.Flags().GetBool("phone-home-enabled")
		req.SetPhoneHomeEnabled(phoneHomeEnabled)
	}

	osItem, resp, err := apiClient.OperatingSystemAPI.CreateOperatingSystem(ctx, org).OperatingSystemCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating operating system (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating operating system: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Operating system created: %s (%s)\n", ptrStr(osItem.Name), ptrStr(osItem.Id))
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	if outputFlag == "yaml" {
		return printYAML(os.Stdout, osItem)
	}
	return printJSON(os.Stdout, osItem)
}

func runOperatingSystemUpdate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewOperatingSystemUpdateRequest()
	changed := false

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		req.SetName(name)
		changed = true
	}
	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		req.SetDescription(description)
		changed = true
	}
	if cmd.Flags().Changed("ipxe-script") {
		ipxeScript, _ := cmd.Flags().GetString("ipxe-script")
		req.SetIpxeScript(ipxeScript)
		changed = true
	}
	if cmd.Flags().Changed("image-url") {
		imageURL, _ := cmd.Flags().GetString("image-url")
		req.SetImageUrl(imageURL)
		changed = true
	}
	if cmd.Flags().Changed("image-sha") {
		imageSHA, _ := cmd.Flags().GetString("image-sha")
		req.SetImageSha(imageSHA)
		changed = true
	}
	if cmd.Flags().Changed("image-auth-type") {
		imageAuthType, _ := cmd.Flags().GetString("image-auth-type")
		req.SetImageAuthType(imageAuthType)
		changed = true
	}
	if cmd.Flags().Changed("image-auth-token") {
		imageAuthToken, _ := cmd.Flags().GetString("image-auth-token")
		req.SetImageAuthToken(imageAuthToken)
		changed = true
	}
	if cmd.Flags().Changed("image-disk") {
		imageDisk, _ := cmd.Flags().GetString("image-disk")
		req.SetImageDisk(imageDisk)
		changed = true
	}
	if cmd.Flags().Changed("root-fs-id") {
		rootFsID, _ := cmd.Flags().GetString("root-fs-id")
		req.SetRootFsId(rootFsID)
		changed = true
	}
	if cmd.Flags().Changed("root-fs-label") {
		rootFsLabel, _ := cmd.Flags().GetString("root-fs-label")
		req.SetRootFsLabel(rootFsLabel)
		changed = true
	}
	if cmd.Flags().Changed("user-data") {
		userData, _ := cmd.Flags().GetString("user-data")
		req.SetUserData(userData)
		changed = true
	}
	if cmd.Flags().Changed("deactivation-note") {
		deactivationNote, _ := cmd.Flags().GetString("deactivation-note")
		req.SetDeactivationNote(deactivationNote)
		changed = true
	}
	if cmd.Flags().Changed("is-cloud-init") {
		isCloudInit, _ := cmd.Flags().GetBool("is-cloud-init")
		req.SetIsCloudInit(isCloudInit)
		changed = true
	}
	if cmd.Flags().Changed("allow-override") {
		allowOverride, _ := cmd.Flags().GetBool("allow-override")
		req.SetAllowOverride(allowOverride)
		changed = true
	}
	if cmd.Flags().Changed("phone-home-enabled") {
		phoneHomeEnabled, _ := cmd.Flags().GetBool("phone-home-enabled")
		req.SetPhoneHomeEnabled(phoneHomeEnabled)
		changed = true
	}
	if cmd.Flags().Changed("is-active") {
		isActive, _ := cmd.Flags().GetBool("is-active")
		req.SetIsActive(isActive)
		changed = true
	}
	if !changed {
		return fmt.Errorf("no update fields provided; specify at least one flag")
	}

	osItem, resp, err := apiClient.OperatingSystemAPI.UpdateOperatingSystem(ctx, org, args[0]).OperatingSystemUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating operating system (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating operating system: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Operating system updated: %s (%s)\n", ptrStr(osItem.Name), ptrStr(osItem.Id))
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	if outputFlag == "yaml" {
		return printYAML(os.Stdout, osItem)
	}
	return printJSON(os.Stdout, osItem)
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
