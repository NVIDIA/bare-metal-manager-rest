// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/nvidia/bare-metal-manager-rest/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ipBlockCmd = &cobra.Command{
	Use:   "ip-block",
	Short: "IP block operations",
}

var ipBlockCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an IP block",
	RunE:  runIPBlockCreate,
}

func init() {
	ipBlockCreateCmd.Flags().String("name", "", "name for the IP block (required)")
	ipBlockCreateCmd.Flags().String("site-id", "", "site ID (required)")
	ipBlockCreateCmd.Flags().String("prefix", "", "IP prefix e.g. 10.0.0.0 (required)")
	ipBlockCreateCmd.Flags().Int32("prefix-length", 0, "prefix length e.g. 16 (required)")
	ipBlockCreateCmd.MarkFlagRequired("name")
	ipBlockCreateCmd.MarkFlagRequired("site-id")
	ipBlockCreateCmd.MarkFlagRequired("prefix")
	ipBlockCreateCmd.MarkFlagRequired("prefix-length")

	rootCmd.AddCommand(ipBlockCmd)
	ipBlockCmd.AddCommand(ipBlockCreateCmd)
}

func runIPBlockCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required")
	}

	name, _ := cmd.Flags().GetString("name")
	siteID, _ := cmd.Flags().GetString("site-id")
	prefix, _ := cmd.Flags().GetString("prefix")
	prefixLen, _ := cmd.Flags().GetInt32("prefix-length")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewIpBlockCreateRequest(name, siteID, "DatacenterOnly", prefix, prefixLen, "IPv4")

	block, resp, err := apiClient.IPBlockAPI.CreateIpblock(ctx, org).IpBlockCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating IP block (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating IP block: %v", err)
	}

	fmt.Fprintf(os.Stderr, "IP block created: %s (%s)\n", ptrStr(block.Name), ptrStr(block.Id))
	return printJSON(os.Stdout, block)
}
