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

// InfiniBandPartitionExpansion - InfiniBandPartition Expansion
type InfiniBandPartitionExpansion interface{}

// InfiniBandPartitionInterface - interface to InfiniBandPartition
type InfiniBandPartitionInterface interface {
	// List all the apis of InfiniBandPartition here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	// Cloud Workflow APIs
	CreateInfiniBandPartition(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.CreateInfiniBandPartitionRequest) (err error)
	DeleteInfiniBandPartition(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.DeleteInfiniBandPartitionRequest) (err error)

	// CreateInfiniBandPartition
	// RegisterWorkflows() error
	// RegisterActivities() error
	GetState() []string
	InfiniBandPartitionExpansion
}
