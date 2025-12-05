// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package networksecuritygroup

import "fmt"

// Init initializes NetworkSecurityGroup Manager
func (api *API) Init() {
	ManagerAccess.Data.EB.Log.Info().Msg("NetworkSecurityGroup: Initializing")
}

// GetState returns the state of NetworkSecurityGroup Manager
func (api *API) GetState() []string {
	state := ManagerAccess.Data.EB.Managers.Workflow.NetworkSecurityGroupState
	var strs []string
	strs = append(strs, fmt.Sprintln("networksecuritygroup_workflow_started", state.WflowStarted.Load()))
	strs = append(strs, fmt.Sprintln("networksecuritygroup_workflow_activity_failed", state.WflowActFail.Load()))
	strs = append(strs, fmt.Sprintln("networksecuritygroup_workflow_activity_succeeded", state.WflowActSucc.Load()))
	strs = append(strs, fmt.Sprintln("networksecuritygroup_workflow_publishing_failed", state.WflowPubFail.Load()))
	strs = append(strs, fmt.Sprintln("networksecuritygroup_workflow_publishing_succeeded", state.WflowPubSucc.Load()))

	return strs
}
