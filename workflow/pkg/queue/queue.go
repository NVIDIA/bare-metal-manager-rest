// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package queue

const (
	// CloudTaskQueue handles all tasks triggered by Cloud API and
	// are meant to be consumed by Cloud system worker
	CloudTaskQueue = "cloud"
	// SiteTaskQueue handles tasks submitted by Site agents running on Site management clusters
	SiteTaskQueue = "site"
)
