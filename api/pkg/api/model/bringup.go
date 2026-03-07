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
	"fmt"

	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
)

// ========== Bring Up Request ==========

// APIBringUpRackRequest is the request body for bring up operations on a single rack
type APIBringUpRackRequest struct {
	SiteID      string `json:"siteId"`
	Description string `json:"description,omitempty"`
}

// Validate validates the bring up request
func (r *APIBringUpRackRequest) Validate() error {
	if r.SiteID == "" {
		return fmt.Errorf("siteId is required")
	}
	return nil
}

// ========== Bring Up Response ==========

// APIBringUpRackResponse is the API response for bring up operations
type APIBringUpRackResponse struct {
	TaskIDs []string `json:"taskIds"`
}

// FromProto converts an RLA SubmitTaskResponse to an APIBringUpRackResponse
func (r *APIBringUpRackResponse) FromProto(resp *rlav1.SubmitTaskResponse) {
	if resp == nil {
		r.TaskIDs = []string{}
		return
	}
	r.TaskIDs = make([]string, 0, len(resp.GetTaskIds()))
	for _, id := range resp.GetTaskIds() {
		r.TaskIDs = append(r.TaskIDs, id.GetId())
	}
}

// NewAPIBringUpRackResponse creates an APIBringUpRackResponse from an RLA SubmitTaskResponse
func NewAPIBringUpRackResponse(resp *rlav1.SubmitTaskResponse) *APIBringUpRackResponse {
	r := &APIBringUpRackResponse{}
	r.FromProto(resp)
	return r
}

// ========== Batch Bring Up Rack Request ==========

// APIBatchBringUpRackRequest is the JSON body for batch rack bring up.
type APIBatchBringUpRackRequest struct {
	SiteID      string      `json:"siteId"`
	Filter      *RackFilter `json:"filter,omitempty"`
	Description string      `json:"description,omitempty"`
}

// Validate checks required fields.
func (r *APIBatchBringUpRackRequest) Validate() error {
	if r.SiteID == "" {
		return fmt.Errorf("siteId is required")
	}
	return nil
}
