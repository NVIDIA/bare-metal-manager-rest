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

// ========== Rack Query Fields ==========

// RackFilterFieldMap maps API field names to RLA protobuf filter enum
var RackFilterFieldMap = map[string]rlav1.RackFilterField{
	"name":         rlav1.RackFilterField_RACK_FILTER_FIELD_NAME,
	"manufacturer": rlav1.RackFilterField_RACK_FILTER_FIELD_MANUFACTURER,
	"model":        rlav1.RackFilterField_RACK_FILTER_FIELD_MODEL,
}

// RackOrderByFieldMap maps API field names to RLA protobuf order by enum
var RackOrderByFieldMap = map[string]rlav1.RackOrderByField{
	"name":         rlav1.RackOrderByField_RACK_ORDER_BY_FIELD_NAME,
	"manufacturer": rlav1.RackOrderByField_RACK_ORDER_BY_FIELD_MANUFACTURER,
	"model":        rlav1.RackOrderByField_RACK_ORDER_BY_FIELD_MODEL,
}

// GetProtoRackFilterFromQueryParam creates an RLA protobuf filter from API query parameters
func GetProtoRackFilterFromQueryParam(fieldName, value string) *rlav1.Filter {
	field, ok := RackFilterFieldMap[fieldName]
	if !ok {
		return nil
	}
	return &rlav1.Filter{
		Field: &rlav1.Filter_RackField{
			RackField: field,
		},
		QueryInfo: &rlav1.StringQueryInfo{
			Patterns:   []string{value},
			IsWildcard: false,
			UseOr:      false,
		},
	}
}

// GetProtoRackOrderByFromQueryParam creates an RLA protobuf OrderBy from API query parameters
func GetProtoRackOrderByFromQueryParam(fieldName, direction string) *rlav1.OrderBy {
	field, ok := RackOrderByFieldMap[fieldName]
	if !ok {
		return nil
	}
	return &rlav1.OrderBy{
		Field: &rlav1.OrderBy_RackField{
			RackField: field,
		},
		Direction: direction,
	}
}

// ========== Rack API Models ==========

// APIRack is the API representation of a Rack from RLA
type APIRack struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Manufacturer string              `json:"manufacturer"`
	Model        string              `json:"model"`
	SerialNumber string              `json:"serialNumber"`
	Description  string              `json:"description"`
	Location     *APIRackLocation    `json:"location,omitempty"`
	Components   []*APIRackComponent `json:"components,omitempty"`
}

// FromProto converts an RLA protobuf Rack to an APIRack
func (ar *APIRack) FromProto(protoRack *rlav1.Rack, includeComponents bool) {
	if protoRack == nil {
		return
	}

	// Get info from DeviceInfo
	if protoRack.GetInfo() != nil {
		info := protoRack.GetInfo()
		if info.GetId() != nil {
			ar.ID = info.GetId().GetId()
		}
		ar.Name = info.GetName()
		ar.Manufacturer = info.GetManufacturer()
		if info.Model != nil {
			ar.Model = *info.Model
		}
		ar.SerialNumber = info.GetSerialNumber()
		if info.Description != nil {
			ar.Description = *info.Description
		}
	}

	// Get location
	if protoRack.GetLocation() != nil {
		ar.Location = &APIRackLocation{}
		ar.Location.FromProto(protoRack.GetLocation())
	}

	// Get components
	if includeComponents && len(protoRack.GetComponents()) > 0 {
		ar.Components = make([]*APIRackComponent, 0, len(protoRack.GetComponents()))
		for _, comp := range protoRack.GetComponents() {
			apiComp := &APIRackComponent{}
			apiComp.FromProto(comp)
			ar.Components = append(ar.Components, apiComp)
		}
	}
}

// NewAPIRack creates an APIRack from the RLA protobuf Rack
func NewAPIRack(protoRack *rlav1.Rack, includeComponents bool) *APIRack {
	if protoRack == nil {
		return nil
	}
	apiRack := &APIRack{}
	apiRack.FromProto(protoRack, includeComponents)
	return apiRack
}

// APIRackLocation represents the location of a rack
type APIRackLocation struct {
	Region     string `json:"region"`
	Datacenter string `json:"datacenter"`
	Room       string `json:"room"`
	Position   string `json:"position"`
}

// FromProto converts a proto Location to an APIRackLocation
func (arl *APIRackLocation) FromProto(protoLocation *rlav1.Location) {
	if protoLocation == nil {
		return
	}
	arl.Region = protoLocation.GetRegion()
	arl.Datacenter = protoLocation.GetDatacenter()
	arl.Room = protoLocation.GetRoom()
	arl.Position = protoLocation.GetPosition()
}

// APIRackComponent represents a component within a rack
type APIRackComponent struct {
	ID              string `json:"id"`
	ComponentID     string `json:"componentId"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	SerialNumber    string `json:"serialNumber"`
	Manufacturer    string `json:"manufacturer"`
	FirmwareVersion string `json:"firmwareVersion"`
	Position        int32  `json:"position"`
}

// FromProto converts a proto Component to an APIRackComponent
func (arc *APIRackComponent) FromProto(protoComponent *rlav1.Component) {
	if protoComponent == nil {
		return
	}
	arc.Type = protoComponent.GetType().String()
	arc.FirmwareVersion = protoComponent.GetFirmwareVersion()
	arc.ComponentID = protoComponent.GetComponentId()

	// Get component info
	if protoComponent.GetInfo() != nil {
		compInfo := protoComponent.GetInfo()
		if compInfo.GetId() != nil {
			arc.ID = compInfo.GetId().GetId()
		}
		arc.Name = compInfo.GetName()
		arc.SerialNumber = compInfo.GetSerialNumber()
		arc.Manufacturer = compInfo.GetManufacturer()
	}

	// Get position
	if protoComponent.GetPosition() != nil {
		arc.Position = protoComponent.GetPosition().GetSlotId()
	}
}
