// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package carbidetypes

import (
	"github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/client"
	"go.uber.org/atomic"
)

// State - Carbide state
type State struct {
	// GrpcFail the number of times the rpc has failed
	GrpcFail atomic.Uint64
	// GrpcSucc the number of times the rpc has succeeded
	GrpcSucc atomic.Uint64
	// HealthStatus current health state
	HealthStatus atomic.Uint64
	// Err is error message
	Err string
	// WflowMetrics workflow metrics
	WflowMetrics WorkflowMetrics
}

// Carbide represents the gRPC client for Carbide and state
type Carbide struct {
	Client *client.CarbideAtomicClient
	State  *State
}

// NewCarbideInstance creates a new instance of Carbide
func NewCarbideInstance() *Carbide {
	carbide := &Carbide{
		State:  &State{},
		Client: client.NewCarbideAtomicClient(&client.CarbideClientConfig{}),
	}

	return carbide
}

// GetClient returns the Carbide client
func (c *Carbide) GetClient() *client.CarbideClient {
	return c.Client.GetClient()
}
