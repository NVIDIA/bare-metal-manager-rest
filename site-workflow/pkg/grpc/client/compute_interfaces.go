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

import wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"

// ComputeGetter is the interface for compute workflows
type ComputeGetter interface {
	Compute() ComputeInterface
}

// ComputeInterface for machine gRPC apis
type ComputeInterface interface {
	MachineInterface
	// Instance Interface
	InstanceInterface
	// SSHKeyGroup Interface
	SSHKeyGroupInterface
	// OperatingSystem Interface
	OperatingSystemInterface
	// Tenant Interface
	TenantInterface
}

type compute struct {
	// carbide client
	carbide wflows.ForgeClient
}

func newCompute(carbide wflows.ForgeClient) *compute {
	return &compute{carbide: carbide}
}
