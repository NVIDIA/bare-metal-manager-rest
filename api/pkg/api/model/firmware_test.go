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
	"testing"

	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewAPIFirmwareUpgradeResponse(t *testing.T) {
	tests := []struct {
		name     string
		resp     *rlav1.SubmitTaskResponse
		expected *APIFirmwareUpgradeResponse
	}{
		{
			name:     "nil response returns empty task IDs",
			resp:     nil,
			expected: &APIFirmwareUpgradeResponse{TaskIDs: []string{}},
		},
		{
			name: "single task ID",
			resp: &rlav1.SubmitTaskResponse{
				TaskIds: []*rlav1.UUID{{Id: "task-1"}},
			},
			expected: &APIFirmwareUpgradeResponse{TaskIDs: []string{"task-1"}},
		},
		{
			name: "multiple task IDs",
			resp: &rlav1.SubmitTaskResponse{
				TaskIds: []*rlav1.UUID{{Id: "task-1"}, {Id: "task-2"}},
			},
			expected: &APIFirmwareUpgradeResponse{TaskIDs: []string{"task-1", "task-2"}},
		},
		{
			name:     "empty task IDs",
			resp:     &rlav1.SubmitTaskResponse{},
			expected: &APIFirmwareUpgradeResponse{TaskIDs: []string{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewAPIFirmwareUpgradeResponse(tt.resp)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAPIRackFirmwareUpgradeBatchRequest_QueryTags(t *testing.T) {
	r := APIRackFirmwareUpgradeBatchRequest{
		SiteID: "site-1",
		Names:  []string{"rack-1", "rack-2"},
	}
	assert.Equal(t, "site-1", r.SiteID)
	assert.Equal(t, []string{"rack-1", "rack-2"}, r.Names)
}

func TestAPIRackFirmwareUpgradeBatchRequest_ToTargetSpec(t *testing.T) {
	tests := []struct {
		name          string
		request       APIRackFirmwareUpgradeBatchRequest
		expectRackLen int
	}{
		{
			name:          "no filter - empty target",
			request:       APIRackFirmwareUpgradeBatchRequest{},
			expectRackLen: 1,
		},
		{
			name: "with name filter",
			request: APIRackFirmwareUpgradeBatchRequest{
				Names: []string{"rack-1", "rack-2"},
			},
			expectRackLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := tt.request.ToTargetSpec()
			assert.NotNil(t, spec)
			racks := spec.GetRacks()
			assert.NotNil(t, racks)
			assert.Len(t, racks.GetTargets(), tt.expectRackLen)
		})
	}
}
