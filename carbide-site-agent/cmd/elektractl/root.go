/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package main

import (
	"fmt"
	"os"

	computils "github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "elektractl",
	Short: "CLI for elektra-site-agent",
	Long:  `CLI for elektra-site-agent`,
}

// statusCmd represents the bootstrap command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "status info",
	Long:  `status info`,
	Run: func(cmd *cobra.Command, args []string) {
		computils.GetSAStatus(computils.SiteStatus)
	},
}

// datastoreCmd represents the datastore command
var datastoreCmd = &cobra.Command{
	Use:   "datastore",
	Short: "datastore cli",
	Long:  `datastore cli`,
	Run: func(cmd *cobra.Command, args []string) {
		computils.GetSAStatus(computils.DatastoreStatus)
	},
}

// machineCmd represents the machine command
var machineCmd = &cobra.Command{
	Use:   "machine",
	Short: "machine cli",
	Long:  `machine cli`,
	Run: func(cmd *cobra.Command, args []string) {
		computils.GetSAStatus(computils.MachineStatus)
	},
}

// vpcCmd represents the vpc command
var vpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "vpc cli",
	Long:  `vpc cli`,
	Run: func(cmd *cobra.Command, args []string) {
		computils.GetSAStatus(computils.VPCStatus)
	},
}

// subnetCmd represents the subnet command
var subnetCmd = &cobra.Command{
	Use:   "subnet",
	Short: "subnet cli",
	Long:  `subnet cli`,
	Run: func(cmd *cobra.Command, args []string) {
		computils.GetSAStatus(computils.SubnetStatus)
	},
}

// instanceCmd represents the instance command
var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "instance cli",
	Long:  `instance cli`,
	Run: func(cmd *cobra.Command, args []string) {
		computils.GetSAStatus(computils.InstanceStatus)
	},
}

// createVpcCmd represents the create command
var createVpcCmd = &cobra.Command{
	Use:   "create",
	Short: "Create VPC",
	Long:  `Create VPC`,
	//Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create Vpc called", vpcName, vpcOrg)
		vpcCreate(vpcName, vpcOrg)
	},
}

// updateCmd represents the update command
var updateVpcCmd = &cobra.Command{
	Use:   "update",
	Short: "Update VPC",
	Long:  `Update VPC`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update Vpc called")
		vpcUpdate(vpcName, vpcOrg, resourceID)
	},
}
var deleteVpcCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete VPC",
	Long:  `Delete VPC`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete Vpc called")
		vpcDelete(resourceID)
	},
}

var getVpcCmd = &cobra.Command{
	Use:   "get",
	Short: "Get VPC",
	Long:  `Get VPC`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Get Vpc called")
		req := computils.VPCStatus + "?" + computils.ParamName + "=" + vpcName
		computils.GetSAStatus(req)
	},
}

// createSubnetCmd represents the create command
var createSubnetCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Subnet",
	Long:  `Create Subnet`,
	Run: func(cmd *cobra.Command, args []string) {
		subnetCreate()
		fmt.Println("create Subnet called ")
	},
}

// updateSubnetCmd represents the update command
var updateSubnetCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Subnet",
	Long:  `Update Subnet`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update Subnet called")
	},
}
var deleteSubnetCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Subnet",
	Long:  `Delete Subnet`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete Subnet called")
	},
}

// createInstanceCmd represents the create command
var createInstanceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Instance",
	Long:  `Create Instance`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create Instance called")
		instanceCreate()
	},
}

// updateInstanceCmd represents the update command
var updateInstanceCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Instance",
	Long:  `Update Instance`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update Instance called")
	},
}
var deleteInstanceCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Instance",
	Long:  `Delete Instance`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete Instance called")
		instanceDelete(resourceID)
	},
}
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Run as cloud proxy forever",
	Long:  `Run as cloud proxy forever`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proxy called")
		proxy()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var vpcName string
var vpcOrg string
var resourceID string

func init() {
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(vpcCmd)
	rootCmd.AddCommand(subnetCmd)
	rootCmd.AddCommand(instanceCmd)
	rootCmd.AddCommand(proxyCmd)
	rootCmd.AddCommand(machineCmd)
	rootCmd.AddCommand(datastoreCmd)

	vpcCmd.AddCommand(createVpcCmd)
	createVpcCmd.Flags().StringVarP(&vpcName, "name", "n", DefaultVPCName, "VPC Name")
	createVpcCmd.Flags().StringVarP(&vpcOrg, "org", "o", DefaultVPCOrg, "VPC Org")
	vpcCmd.AddCommand(updateVpcCmd)
	updateVpcCmd.Flags().StringVarP(&vpcName, "name", "n", DefaultVPCName, "VPC Name")
	updateVpcCmd.Flags().StringVarP(&vpcOrg, "org", "o", DefaultVPCOrg, "VPC Org")
	vpcCmd.AddCommand(deleteVpcCmd)
	deleteVpcCmd.Flags().StringVarP(&resourceID, "id", "i", "", "VPC Id")
	deleteVpcCmd.MarkFlagRequired("id")
	vpcCmd.AddCommand(getVpcCmd)
	getVpcCmd.Flags().StringVarP(&vpcName, "name", "n", DefaultVPCName, "VPC Name")

	subnetCmd.AddCommand(createSubnetCmd)
	subnetCmd.AddCommand(updateSubnetCmd)
	subnetCmd.AddCommand(deleteSubnetCmd)

	instanceCmd.AddCommand(createInstanceCmd)
	instanceCmd.AddCommand(updateInstanceCmd)
	instanceCmd.AddCommand(deleteInstanceCmd)
	deleteInstanceCmd.Flags().StringVarP(&resourceID, "id", "i", "", "Instance Id")
}
