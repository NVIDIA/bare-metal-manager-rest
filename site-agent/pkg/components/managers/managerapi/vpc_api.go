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

// VPCExpansion - VPC Expansion
type VPCExpansion interface{}

// VPCInterface - interface to VPC
type VPCInterface interface {
	// List all the apis of VPC here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	// Cloud Workflow APIs
	CreateVPC(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.CreateVPCRequest) (err error)
	DeleteVPC(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.DeleteVPCRequest) (err error)
	// 	UpdateVpcInfo(ctx workflow.Context, SiteID string, TransactionID *wflows.TransactionID, VPCInfo *wflows.VPCInfo) (err error)

	// CRUD VPC APIs
	UpdateVPC(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.UpdateVPCRequest) (err error)
	// GetVPCByID(ctx workflow.Context, ResourceID string, VPCID string) (ResourceResponse *wflows.GetVPCResponse, err error)
	GetVPCByName(ctx workflow.Context, ResourceID string, VPCName string) (ResourceResponse *wflows.GetVPCResponse, err error)
	// GetVPCAll(ctx workflow.Context, ResourceID string) (ResourceResponse *wflows.GetVPCResponse, err error)
	// DeleteVPCByIDWorkflow(ctx workflow.Context, ResourceID string, VPCID string) (err error)

	// CreateVPC
	// RegisterWorkflows() error
	// RegisterActivities() error
	GetState() []string
	VPCExpansion
}
