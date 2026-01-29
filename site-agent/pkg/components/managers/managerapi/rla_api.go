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
	"context"

	"github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/client"
)

// RLAExpansion - RLA Expansion
type RLAExpansion interface{}

// RLAInterface - interface to RLA
type RLAInterface interface {
	// List all the apis of RLA here
	Init()
	Start()
	CreateGrpcClient() error
	GetGrpcClient() *client.RlaClient
	UpdateGrpcClientState(err error)
	CreateGrpcClientActivity(ctx context.Context, ResourceID string) (client *client.RlaClient, err error)
	RegisterGrpc()
	RegisterSubscriber() error
	GetState() []string
	GetGrpcClientVersion() int64
	RLAExpansion
}
