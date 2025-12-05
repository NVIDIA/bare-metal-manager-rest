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

// NVExpansion - ExpectedMachine Expansion
type NVLinkLogicalPartitionExpansion interface{}

// NVLinkLogicalPartitionInterface - interface to NVLinkLogicalPartition
type NVLinkLogicalPartitionInterface interface {
	// List all the apis of NVLinkLogicalPartition here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	GetState() []string
	NVLinkLogicalPartitionExpansion
}
