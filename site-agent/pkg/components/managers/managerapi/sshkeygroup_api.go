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

// SSHKeyGroupExpansion - SSHKeyGroup Expansion
type SSHKeyGroupExpansion interface{}

// SSHKeyGroupInterface - interface to SSHKeyGroup
type SSHKeyGroupInterface interface {
	// List all the apis of SSHKeyGroup here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	// Cloud Workflow APIs
	CreateSSHKeyGroup(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.CreateSSHKeyGroupRequest) (err error)
	DeleteSSHKeyGroup(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.DeleteSSHKeyGroupRequest) (err error)

	// CRUD SSHKeyGroup APIs
	UpdateSSHKeyGroup(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.UpdateSSHKeyGroupRequest) (err error)
	// GetSSHKeyGroupByID(ctx workflow.Context, ResourceID string, SSHKeyGroupID string) (ResourceResponse *wflows.GetSSHKeyGroupResponse, err error)
	GetSSHKeyGroup(ctx workflow.Context, ResourceID string, ResourceRequest *wflows.GetSSHKeyGroup) (ResourceResponse *wflows.GetSSHKeyGroupResponse, err error)

	// CreateSSHKeyGroup
	// RegisterWorkflows() error
	// RegisterActivities() error
	GetState() []string
	SSHKeyGroupExpansion
}
