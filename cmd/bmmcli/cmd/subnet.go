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

var subnetCmd = &cobra.Command{
	Use:   "subnet",
	Short: "Subnet operations",
}

var subnetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a subnet",
	RunE:  runSubnetCreate,
}

func init() {
	subnetCreateCmd.Flags().String("name", "", "name for the subnet (required)")
	subnetCreateCmd.Flags().String("vpc-id", "", "VPC ID (required)")
	subnetCreateCmd.Flags().String("ipv4-block-id", "", "IPv4 block ID (required)")
	subnetCreateCmd.Flags().Int32("prefix-length", 24, "prefix length (default 24)")
	subnetCreateCmd.MarkFlagRequired("name")
	subnetCreateCmd.MarkFlagRequired("vpc-id")
	subnetCreateCmd.MarkFlagRequired("ipv4-block-id")

	rootCmd.AddCommand(subnetCmd)
	subnetCmd.AddCommand(subnetCreateCmd)
}

func runSubnetCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required")
	}

	name, _ := cmd.Flags().GetString("name")
	vpcID, _ := cmd.Flags().GetString("vpc-id")
	blockID, _ := cmd.Flags().GetString("ipv4-block-id")
	prefixLen, _ := cmd.Flags().GetInt32("prefix-length")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewSubnetCreateRequest(name, vpcID, prefixLen)
	req.Ipv4BlockId = &blockID

	subnet, resp, err := apiClient.SubnetAPI.CreateSubnet(ctx, org).SubnetCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating subnet (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating subnet: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Subnet created: %s (%s)\n", ptrStr(subnet.Name), ptrStr(subnet.Id))
	return printJSON(os.Stdout, subnet)
}
