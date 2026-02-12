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

var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Instance operations",
}

var instanceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an instance",
	RunE:  runInstanceCreate,
}

func init() {
	instanceCreateCmd.Flags().String("name", "", "name for the instance (required)")
	instanceCreateCmd.Flags().String("vpc-id", "", "VPC ID (required)")
	instanceCreateCmd.Flags().String("subnet-id", "", "subnet ID for the network interface (required)")
	instanceCreateCmd.Flags().String("instance-type-id", "", "instance type ID (required)")
	instanceCreateCmd.Flags().String("tenant-id", "", "tenant ID (required)")
	instanceCreateCmd.Flags().String("os-id", "", "operating system ID (optional)")
	instanceCreateCmd.Flags().String("ipxe-script", "", "iPXE script (optional, used if no OS)")
	instanceCreateCmd.MarkFlagRequired("name")
	instanceCreateCmd.MarkFlagRequired("vpc-id")
	instanceCreateCmd.MarkFlagRequired("subnet-id")
	instanceCreateCmd.MarkFlagRequired("instance-type-id")
	instanceCreateCmd.MarkFlagRequired("tenant-id")

	rootCmd.AddCommand(instanceCmd)
	instanceCmd.AddCommand(instanceCreateCmd)
}

func runInstanceCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required")
	}

	name, _ := cmd.Flags().GetString("name")
	vpcID, _ := cmd.Flags().GetString("vpc-id")
	subnetID, _ := cmd.Flags().GetString("subnet-id")
	itID, _ := cmd.Flags().GetString("instance-type-id")
	tenantID, _ := cmd.Flags().GetString("tenant-id")
	osID, _ := cmd.Flags().GetString("os-id")
	ipxeScript, _ := cmd.Flags().GetString("ipxe-script")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	iface := client.InterfaceCreateRequest{
		SubnetId: &subnetID,
	}

	req := client.NewInstanceCreateRequest(name, tenantID, vpcID, []client.InterfaceCreateRequest{iface})
	req.InstanceTypeId = &itID

	if osID != "" {
		req.SetOperatingSystemId(osID)
	} else if ipxeScript != "" {
		req.SetIpxeScript(ipxeScript)
	} else {
		req.SetIpxeScript("#!ipxe\necho No OS configured")
	}

	instance, resp, err := apiClient.InstanceAPI.CreateInstance(ctx, org).InstanceCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating instance (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating instance: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Instance created: %s (%s)\n", ptrStr(instance.Name), ptrStr(instance.Id))
	return printJSON(os.Stdout, instance)
}
