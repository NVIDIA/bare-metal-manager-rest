/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package model

import (
	rlav1 "github.com/nvidia/carbide-rest/workflow-schema/rla/protobuf/v1"
)

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

// NewAPIRack creates an APIRack from the RLA protobuf Rack
func NewAPIRack(rack *rlav1.Rack, withComponents bool) *APIRack {
	if rack == nil {
		return nil
	}

	apiRack := &APIRack{}

	// Get info from DeviceInfo
	if rack.GetInfo() != nil {
		info := rack.GetInfo()
		if info.GetId() != nil {
			apiRack.ID = info.GetId().GetId()
		}
		apiRack.Name = info.GetName()
		apiRack.Manufacturer = info.GetManufacturer()
		if info.Model != nil {
			apiRack.Model = *info.Model
		}
		apiRack.SerialNumber = info.GetSerialNumber()
		if info.Description != nil {
			apiRack.Description = *info.Description
		}
	}

	// Get location
	if rack.GetLocation() != nil {
		apiRack.Location = &APIRackLocation{}
		apiRack.Location.FromProto(rack.GetLocation())
	}

	// Get components
	if withComponents && len(rack.GetComponents()) > 0 {
		apiRack.Components = make([]*APIRackComponent, 0, len(rack.GetComponents()))
		for _, comp := range rack.GetComponents() {
			apiComp := &APIRackComponent{}
			apiComp.FromProto(comp)
			apiRack.Components = append(apiRack.Components, apiComp)
		}
	}

	return apiRack
}
