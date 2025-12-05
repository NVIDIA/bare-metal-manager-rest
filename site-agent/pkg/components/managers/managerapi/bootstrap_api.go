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

// BootstrapExpansion - Bootstrap Expansion
type BootstrapExpansion interface{}

// BootstrapInterface - interface to Bootstrap
type BootstrapInterface interface {
	// List all the apis of Bootstrap here
	Init()
	Start()
	DownloadAndStoreCreds(otpOverride []byte) error
	GetState() []string

	BootstrapExpansion

	// Temporal Workflows - Subscriber
	RegisterSubscriber() error
}
