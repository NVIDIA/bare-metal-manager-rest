/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
)

// ========== Firmware Upgrade Request ==========

// APIFirmwareUpgradeRequest is the request body for firmware upgrade operations
type APIFirmwareUpgradeRequest struct {
	Version *string `json:"version,omitempty"`
}

// ========== Firmware Upgrade Response ==========

// APIFirmwareUpgradeResponse is the API response for firmware upgrade operations
type APIFirmwareUpgradeResponse struct {
	TaskIDs []string `json:"taskIds"`
}

// NewAPIFirmwareUpgradeResponse creates an APIFirmwareUpgradeResponse from an RLA SubmitTaskResponse
func NewAPIFirmwareUpgradeResponse(resp *rlav1.SubmitTaskResponse) *APIFirmwareUpgradeResponse {
	if resp == nil {
		return &APIFirmwareUpgradeResponse{TaskIDs: []string{}}
	}
	taskIDs := make([]string, 0, len(resp.GetTaskIds()))
	for _, id := range resp.GetTaskIds() {
		taskIDs = append(taskIDs, id.GetId())
	}
	return &APIFirmwareUpgradeResponse{TaskIDs: taskIDs}
}

// ========== Rack Firmware Upgrade Batch Request ==========

// APIRackFirmwareUpgradeBatchRequest captures query parameters for batch rack firmware upgrade.
// Supports filtering by rack name.
type APIRackFirmwareUpgradeBatchRequest struct {
	SiteID string   `query:"siteId"`
	Names  []string `query:"name"`
}

// ToTargetSpec converts the filter request to an RLA OperationTargetSpec
func (r *APIRackFirmwareUpgradeBatchRequest) ToTargetSpec() *rlav1.OperationTargetSpec {
	var rackTargets []*rlav1.RackTarget

	for _, name := range r.Names {
		rackTargets = append(rackTargets, &rlav1.RackTarget{
			Identifier: &rlav1.RackTarget_Name{
				Name: name,
			},
		})
	}

	if len(rackTargets) == 0 {
		rackTargets = append(rackTargets, &rlav1.RackTarget{})
	}

	return &rlav1.OperationTargetSpec{
		Targets: &rlav1.OperationTargetSpec_Racks{
			Racks: &rlav1.RackTargets{
				Targets: rackTargets,
			},
		},
	}
}
