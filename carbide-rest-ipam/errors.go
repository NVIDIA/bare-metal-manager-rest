// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package ipam

import "errors"

var (
	// ErrNotFound is returned if prefix or cidr was not found
	ErrNotFound = errors.New("NotFound")
	// ErrNoIPAvailable is returned if no IP is available anymore
	ErrNoIPAvailable = errors.New("NoIPAvailableError")
	// ErrAlreadyAllocated is returned if the requested address is not available
	ErrAlreadyAllocated = errors.New("AlreadyAllocatedError")
	// ErrOptimisticLockError is returned if insert or update conflicts with the existing data
	ErrOptimisticLockError = errors.New("OptimisticLockError")
	// ErrNamespaceDoesNotExist is returned when an operation is perfomed in a namespace that does not exist.
	ErrNamespaceDoesNotExist = errors.New("NamespaceDoesNotExist")
	// ErrNameTooLong is returned when a name exceeds the databases max identifier length
	ErrNameTooLong = errors.New("NameTooLong")
	// ErrNotImplemented is returned for CreateNamespace, ListNamespaces, DeleteNamespace apis
	ErrNotImplemented = errors.New("NotImplemented")
)
