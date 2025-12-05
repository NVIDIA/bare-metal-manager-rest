// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package managerapi

import (
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"

	"go.temporal.io/sdk/workflow"
)

// SubnetExpansion - Subnet Expansion
type SubnetExpansion interface{}

// SubnetInterface - interface to Subnet
type SubnetInterface interface {
	// List all the apis of Subnet here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	// Temporal Workflows - Subscriber
	CreateSubnet(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.CreateSubnetRequest) (err error)
	// Implement this when this is available in Site controller
	// UpdateSubnet(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.UpdateSubnetRequest) (err error)
	DeleteSubnet(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.DeleteSubnetRequest) (err error)
	GetState() []string
	SubnetExpansion
}
