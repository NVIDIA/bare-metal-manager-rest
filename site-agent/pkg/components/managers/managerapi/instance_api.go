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

// InstanceExpansion - Instance Expansion
type InstanceExpansion interface{}

// InstanceInterface - interface to Instance
type InstanceInterface interface {
	// List all the apis of Instance here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	// Temporal Workflows - Subscriber
	//Create Instance (deprecated)
	CreateInstance(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.CreateInstanceRequest) (err error)
	//Delete Instance (deprecated)
	DeleteInstance(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.DeleteInstanceRequest) (err error)
	//RebootInstance (deprecated)
	RebootInstance(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.RebootInstanceRequest) (err error)
	GetState() []string

	InstanceExpansion
}
