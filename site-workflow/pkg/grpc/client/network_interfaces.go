// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package client

import (
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

// NetworkGetter is the interface for the network workflows
type NetworkGetter interface {
	Network() NetworkInterface
}

// NetworkInterface is the interface for the network client
type NetworkInterface interface {

	// VPC Interface
	VPCInterface
	// Subnet Interface
	SubnetInterface
	// InfiniBandPartition Interface
	InfiniBandPartitionInterface
}

type network struct {
	// carbide client
	carbide wflows.ForgeClient
}

func newNetwork(carbide wflows.ForgeClient) *network {
	return &network{carbide: carbide}
}
