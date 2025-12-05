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

// CarbideExpansion - Carbide Expansion
type CarbideExpansion interface{}

// CarbideInterface - interface to Carbide
type CarbideInterface interface {
	// List all the apis of Carbide here
	Init()
	Start()
	CreateGRPCClient() error
	GetGRPCClient() *client.CarbideClient
	UpdateGRPCClientState(err error)
	CreateGRPCClientActivity(ctx context.Context, ResourceID string) (client *client.CarbideClient, err error)
	RegisterGRPC()
	GetState() []string
	GetGRPCClientVersion() int64
	CarbideExpansion
}
