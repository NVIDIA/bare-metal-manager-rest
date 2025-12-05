// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package instancetype

import "fmt"

// Init initializes Instance Type Manager
func (api *API) Init() {
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: Initializing the Subnet")
}

// GetState returns the state of Instance Type Manager
func (api *API) GetState() []string {
	state := ManagerAccess.Data.EB.Managers.Workflow.InstanceTypeState
	var strs []string
	strs = append(strs, fmt.Sprintln("instancetype_workflow_started", state.WflowStarted.Load()))
	strs = append(strs, fmt.Sprintln("instancetype_workflow_activity_failed", state.WflowActFail.Load()))
	strs = append(strs, fmt.Sprintln("instancetype_worflow_activity_succeeded", state.WflowActSucc.Load()))
	strs = append(strs, fmt.Sprintln("instancetype_workflow_publishing_failed", state.WflowPubFail.Load()))
	strs = append(strs, fmt.Sprintln("instancetype_worflow_publishing_succeeded", state.WflowPubSucc.Load()))

	return strs
}
