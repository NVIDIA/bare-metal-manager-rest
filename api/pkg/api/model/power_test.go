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

func TestAPIPowerControlRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request APIPowerControlRequest
		wantErr bool
	}{
		{
			name:    "valid - on",
			request: APIPowerControlRequest{State: "on"},
			wantErr: false,
		},
		{
			name:    "valid - off",
			request: APIPowerControlRequest{State: "off"},
			wantErr: false,
		},
		{
			name:    "valid - cycle",
			request: APIPowerControlRequest{State: "cycle"},
			wantErr: false,
		},
		{
			name:    "valid - forceoff",
			request: APIPowerControlRequest{State: "forceoff"},
			wantErr: false,
		},
		{
			name:    "valid - forcecycle",
			request: APIPowerControlRequest{State: "forcecycle"},
			wantErr: false,
		},
		{
			name:    "invalid - empty state",
			request: APIPowerControlRequest{State: ""},
			wantErr: true,
		},
		{
			name:    "invalid - unknown state",
			request: APIPowerControlRequest{State: "reboot"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewAPIPowerControlResponse(t *testing.T) {
	tests := []struct {
		name     string
		resp     *rlav1.SubmitTaskResponse
		expected *APIPowerControlResponse
	}{
		{
			name:     "nil response returns empty task IDs",
			resp:     nil,
			expected: &APIPowerControlResponse{TaskIDs: []string{}},
		},
		{
			name: "response with task IDs",
			resp: &rlav1.SubmitTaskResponse{
				TaskIds: []*rlav1.UUID{
					{Id: "task-1"},
					{Id: "task-2"},
				},
			},
			expected: &APIPowerControlResponse{TaskIDs: []string{"task-1", "task-2"}},
		},
		{
			name: "response with empty task IDs",
			resp: &rlav1.SubmitTaskResponse{
				TaskIds: []*rlav1.UUID{},
			},
			expected: &APIPowerControlResponse{TaskIDs: []string{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAPIPowerControlResponse(tt.resp)
			assert.NotNil(t, got)
			assert.Equal(t, tt.expected.TaskIDs, got.TaskIDs)
		})
	}
}

func TestAPIBatchRackPowerControlRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request APIBatchRackPowerControlRequest
		wantErr bool
	}{
		{
			name:    "valid - on with siteId",
			request: APIBatchRackPowerControlRequest{SiteID: "site-1", State: "on"},
			wantErr: false,
		},
		{
			name: "valid - with filter",
			request: APIBatchRackPowerControlRequest{
				SiteID: "site-1",
				Filter: &RackFilter{Names: []string{"Rack-001"}},
				State:  "off",
			},
			wantErr: false,
		},
		{
			name:    "invalid - missing siteId",
			request: APIBatchRackPowerControlRequest{State: "on"},
			wantErr: true,
		},
		{
			name:    "invalid - bad state",
			request: APIBatchRackPowerControlRequest{SiteID: "site-1", State: "reboot"},
			wantErr: true,
		},
		{
			name:    "invalid - empty state",
			request: APIBatchRackPowerControlRequest{SiteID: "site-1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRackFilter_ToTargetSpec(t *testing.T) {
	tests := []struct {
		name          string
		filter        *RackFilter
		expectedRacks int
	}{
		{
			name:          "nil filter - targets all racks",
			filter:        nil,
			expectedRacks: 1,
		},
		{
			name:          "empty filter - targets all racks",
			filter:        &RackFilter{},
			expectedRacks: 1,
		},
		{
			name:          "with single name",
			filter:        &RackFilter{Names: []string{"Rack-001"}},
			expectedRacks: 1,
		},
		{
			name:          "with multiple names",
			filter:        &RackFilter{Names: []string{"Rack-001", "Rack-002"}},
			expectedRacks: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := tt.filter.ToTargetSpec()
			assert.NotNil(t, spec)

			racks := spec.GetRacks()
			assert.NotNil(t, racks)
			assert.Len(t, racks.GetTargets(), tt.expectedRacks)
		})
	}
}
