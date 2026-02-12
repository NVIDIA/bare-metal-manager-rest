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

var instanceTypeCmd = &cobra.Command{
	Use:   "instance-type",
	Short: "Instance type operations",
}

var instanceTypeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an instance type",
	RunE:  runInstanceTypeCreate,
}

func init() {
	instanceTypeCreateCmd.Flags().String("name", "", "name for the instance type (required)")
	instanceTypeCreateCmd.Flags().String("site-id", "", "site ID (required)")
	instanceTypeCreateCmd.Flags().String("description", "", "optional description")
	instanceTypeCreateCmd.MarkFlagRequired("name")
	instanceTypeCreateCmd.MarkFlagRequired("site-id")

	rootCmd.AddCommand(instanceTypeCmd)
	instanceTypeCmd.AddCommand(instanceTypeCreateCmd)
}

func runInstanceTypeCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required")
	}

	name, _ := cmd.Flags().GetString("name")
	siteID, _ := cmd.Flags().GetString("site-id")
	description, _ := cmd.Flags().GetString("description")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewInstanceTypeCreateRequest(name, siteID)
	if description != "" {
		req.Description = &description
	}

	it, resp, err := apiClient.InstanceTypeAPI.CreateInstanceType(ctx, org).InstanceTypeCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating instance type (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating instance type: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Instance type created: %s (%s)\n", ptrStr(it.Name), ptrStr(it.Id))
	return printJSON(os.Stdout, it)
}
