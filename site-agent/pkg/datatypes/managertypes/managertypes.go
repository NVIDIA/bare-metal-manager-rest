// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package managertypes

import (
	bootstraptypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/bootstrap"
	carbidetypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/carbide"
	healthtypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/health"
	workflowtypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/workflow"
)

// Managers - manager ds
type Managers struct {
	Version string
	// All the datastructures of Managers below
	Workflow  *workflowtypes.Workflow
	Carbide   *carbidetypes.Carbide
	Bootstrap *bootstraptypes.Bootstrap
	Health    *healthtypes.HealthCache
}

// NewManagerType - get new type of all managers
func NewManagerType() *Managers {
	return &Managers{
		Version: "0.0.1",
		// All the managers below
		Workflow:  workflowtypes.NewWorkflowInstance(),
		Carbide:   carbidetypes.NewCarbideInstance(),
		Bootstrap: bootstraptypes.NewBootstrapInstance(),
		Health:    healthtypes.NewHealthCache(),
	}
}
