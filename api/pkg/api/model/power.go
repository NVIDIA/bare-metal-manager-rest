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
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
)

// ValidPowerControlStates defines the valid states for power control operations
var ValidPowerControlStates = []string{"on", "off", "cycle", "forceoff", "forcecycle"}

var validPowerControlStatesAny = func() []interface{} {
	result := make([]interface{}, len(ValidPowerControlStates))
	for i, s := range ValidPowerControlStates {
		result[i] = s
	}
	return result
}()

// ========== Power Control Request ==========

// APIPowerControlRequest is the request body for power control operations
type APIPowerControlRequest struct {
	State string `json:"state"`
}

// Validate validates the power control request
func (r *APIPowerControlRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.State,
			validation.Required.Error(validationErrorValueRequired),
			validation.In(validPowerControlStatesAny...).Error(
				fmt.Sprintf("must be one of %v", ValidPowerControlStates))),
	)
}

// ========== Power Control Response ==========

// APIPowerControlResponse is the API response for power control operations
type APIPowerControlResponse struct {
	TaskIDs []string `json:"taskIds"`
}

// FromProto converts an RLA SubmitTaskResponse to an APIPowerControlResponse
func (r *APIPowerControlResponse) FromProto(resp *rlav1.SubmitTaskResponse) {
	if resp == nil {
		r.TaskIDs = []string{}
		return
	}
	r.TaskIDs = make([]string, 0, len(resp.GetTaskIds()))
	for _, id := range resp.GetTaskIds() {
		r.TaskIDs = append(r.TaskIDs, id.GetId())
	}
}

// NewAPIPowerControlResponse creates an APIPowerControlResponse from an RLA SubmitTaskResponse
func NewAPIPowerControlResponse(resp *rlav1.SubmitTaskResponse) *APIPowerControlResponse {
	r := &APIPowerControlResponse{}
	r.FromProto(resp)
	return r
}

// ========== Rack Power Control Batch Request ==========

// APIRackPowerControlBatchRequest captures query parameters for batch rack power control.
// Supports filtering by rack name.
type APIRackPowerControlBatchRequest struct {
	SiteID string   `query:"siteId"`
	Names  []string `query:"name"`
}

// Validate checks required fields on the batch power control request
func (r *APIRackPowerControlBatchRequest) Validate() error {
	if r.SiteID == "" {
		return fmt.Errorf("siteId query parameter is required")
	}
	return nil
}

// QueryValues returns only the known query parameters as url.Values,
// suitable for deterministic workflow ID hashing without unknown param interference.
func (r *APIRackPowerControlBatchRequest) QueryValues() url.Values {
	v := url.Values{}
	v.Set("siteId", r.SiteID)
	for _, n := range r.Names {
		v.Add("name", n)
	}
	return v
}

// ToTargetSpec converts the filter request to an RLA OperationTargetSpec
func (r *APIRackPowerControlBatchRequest) ToTargetSpec() *rlav1.OperationTargetSpec {
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
