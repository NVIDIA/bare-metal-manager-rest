// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package wflowinterface

import (
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"go.temporal.io/sdk/workflow"
)

// VPCInterface - cloud workflow interface for vpc updates
type VPCInterface interface {
	UpdateVpcInfo(ctx workflow.Context, SiteID string, TransactionID *wflows.TransactionID, VPCInfo *wflows.VPCInfo) (err error)
	UpdateVpcInventory(ctx workflow.Context, SiteID string, VPCInventory *wflows.VPCInventory) (err error)
}

// MachineInterface - cloud workflow interface for machine updates
type MachineInterface interface {
	UpdateMachineInventory(ctx workflow.Context, SiteID string, MachineInventory *wflows.MachineInventory) (err error)
}

// SubnetInterface - cloud workflow interface for Subnet updates
type SubnetInterface interface {
	UpdateSubnetInfo(ctx workflow.Context, SiteID string, TransactionID *wflows.TransactionID, SubnetInfo *wflows.SubnetInfo) (err error)
	UpdateSubnetInventory(ctx workflow.Context, SiteID string, SubnetInventory *wflows.SubnetInventory) (err error)
}

// InstanceInterface - cloud workflow interface for Instance updates
type InstanceInterface interface {
	UpdateInstanceInfo(ctx workflow.Context, SiteID string, TransactionID *wflows.TransactionID, Instance *wflows.InstanceInfo) (err error)
	UpdateInstanceRebootInfo(ctx workflow.Context, SiteID string, TransactionID *wflows.TransactionID, Instance *wflows.InstanceRebootInfo) (err error)
	UpdateInstanceInventory(ctx workflow.Context, SiteID string, InstanceInventory *wflows.InstanceInventory) (err error)
}

// SSHKeyGroupInterface - cloud workflow interface for SSHKeyGroup updates
type SSHKeyGroupInterface interface {
	UpdateSSHKeyGroupInfo(ctx workflow.Context, SiteID string, TransactionID *wflows.TransactionID, SSHKeyGroupInfo *wflows.SSHKeyGroupInfo) (err error)
	UpdateSSHKeyGroupInventory(ctx workflow.Context, SiteID string, SSHKeyGroupInventory *wflows.SSHKeyGroupInventory) (err error)
}

// InfiniBandPartitionInterface - cloud workflow interface for InfiniBandPartition updates
type InfiniBandPartitionInterface interface {
	UpdateInfiniBandPartitionInfo(ctx workflow.Context, SiteID string, TransactionID *wflows.TransactionID, InfiniBandPartitionInfo *wflows.InfiniBandPartitionInfo) (err error)
	UpdateInfiniBandPartitionInventory(ctx workflow.Context, SiteID string, InfiniBandPartitionInventory *wflows.InfiniBandPartitionInventory) (err error)
}
