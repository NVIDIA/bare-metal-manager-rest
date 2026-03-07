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

func TestAPIBringUpRackRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request APIBringUpRackRequest
		wantErr bool
	}{
		{
			name:    "valid - with siteId",
			request: APIBringUpRackRequest{SiteID: "site-1"},
			wantErr: false,
		},
		{
			name:    "valid - with siteId and description",
			request: APIBringUpRackRequest{SiteID: "site-1", Description: "bring up rack"},
			wantErr: false,
		},
		{
			name:    "invalid - missing siteId",
			request: APIBringUpRackRequest{},
			wantErr: true,
		},
		{
			name:    "invalid - empty siteId",
			request: APIBringUpRackRequest{SiteID: ""},
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

func TestNewAPIBringUpRackResponse(t *testing.T) {
	tests := []struct {
		name     string
		resp     *rlav1.SubmitTaskResponse
		expected *APIBringUpRackResponse
	}{
		{
			name:     "nil response returns empty task IDs",
			resp:     nil,
			expected: &APIBringUpRackResponse{TaskIDs: []string{}},
		},
		{
			name: "single task ID",
			resp: &rlav1.SubmitTaskResponse{
				TaskIds: []*rlav1.UUID{{Id: "task-1"}},
			},
			expected: &APIBringUpRackResponse{TaskIDs: []string{"task-1"}},
		},
		{
			name: "multiple task IDs",
			resp: &rlav1.SubmitTaskResponse{
				TaskIds: []*rlav1.UUID{{Id: "task-1"}, {Id: "task-2"}},
			},
			expected: &APIBringUpRackResponse{TaskIDs: []string{"task-1", "task-2"}},
		},
		{
			name:     "empty task IDs",
			resp:     &rlav1.SubmitTaskResponse{},
			expected: &APIBringUpRackResponse{TaskIDs: []string{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewAPIBringUpRackResponse(tt.resp)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected.TaskIDs, result.TaskIDs)
		})
	}
}

func TestAPIBatchBringUpRackRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request APIBatchBringUpRackRequest
		wantErr bool
	}{
		{
			name:    "valid - with siteId only",
			request: APIBatchBringUpRackRequest{SiteID: "site-1"},
			wantErr: false,
		},
		{
			name: "valid - with filter",
			request: APIBatchBringUpRackRequest{
				SiteID: "site-1",
				Filter: &RackFilter{Names: []string{"Rack-001"}},
			},
			wantErr: false,
		},
		{
			name: "valid - with description",
			request: APIBatchBringUpRackRequest{
				SiteID:      "site-1",
				Description: "batch bring up",
			},
			wantErr: false,
		},
		{
			name:    "invalid - missing siteId",
			request: APIBatchBringUpRackRequest{},
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
