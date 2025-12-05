// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package healthtypes

import (
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

// We will define our own state later
// type HealthState int

// const (
// 	UNKNOWN HealthState = iota
// 	UP
// 	DOWN
// 	ERROR
// )

type SiteInventoryHealth struct {
	State     wflows.HealthState
	StatusMsg string
	// More fields to be added later
}

type SiteControllerConnection struct {
	State     wflows.HealthState
	StatusMsg string
	// More fields to be added later
}

type HighAvailability struct {
	State     wflows.HealthState
	StatusMsg string
}

// HealthCache Site Agent HealthCache
type HealthCache struct {
	Inventory        SiteInventoryHealth
	CarbideInterface SiteControllerConnection
	Availabilty      HighAvailability
}

// NewHealthCache - initialize site agent health cache
func NewHealthCache() *HealthCache {
	return &HealthCache{}
}
