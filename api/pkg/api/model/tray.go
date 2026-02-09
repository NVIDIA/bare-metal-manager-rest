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
	"strings"

	rlav1 "github.com/nvidia/carbide-rest/workflow-schema/rla/protobuf/v1"
)

// ValidTrayTypes contains the valid tray type strings
var ValidTrayTypes = []string{"compute", "switch", "powershelf", "torswitch", "ums", "cdu"}

// TrayFilterInput is the filter structure for querying trays from RLA
type TrayFilterInput struct {
	// RackID filters trays by rack UUID
	RackID *string
	// RackName filters trays by rack name
	RackName *string
	// Slot filters trays by slot number in the rack
	Slot *int32
	// Type filters trays by type (compute, switch, powershelf, torswitch, ums, cdu)
	Type *string
	// Index filters trays by index of tray in its type
	Index *int32
	// ComponentIDs filters trays by component IDs (comma-separated)
	ComponentIDs []string
	// IDs filters trays by UUIDs (comma-separated)
	IDs []string
	// TaskID filters trays by task UUID
	TaskID *string
}

// APITray is the API representation of a Tray (Component) from RLA
type APITray struct {
	ID              string           `json:"id"`
	ComponentID     string           `json:"componentId"`
	Type            string           `json:"type"`
	Name            string           `json:"name"`
	Manufacturer    string           `json:"manufacturer"`
	Model           string           `json:"model,omitempty"`
	SerialNumber    string           `json:"serialNumber"`
	Description     string           `json:"description,omitempty"`
	FirmwareVersion string           `json:"firmwareVersion"`
	Position        *APITrayPosition `json:"position,omitempty"`
	RackID          string           `json:"rackId,omitempty"`
}

// APITrayPosition represents the position of a tray within a rack
type APITrayPosition struct {
	SlotID  int32 `json:"slotId"`
	TrayIdx int32 `json:"trayIdx"`
	HostID  int32 `json:"hostId"`
}

// componentTypeToString converts a ComponentType enum to a string
func componentTypeToString(ct rlav1.ComponentType) string {
	switch ct {
	case rlav1.ComponentType_COMPONENT_TYPE_COMPUTE:
		return "compute"
	case rlav1.ComponentType_COMPONENT_TYPE_NVLSWITCH:
		return "switch"
	case rlav1.ComponentType_COMPONENT_TYPE_POWERSHELF:
		return "powershelf"
	case rlav1.ComponentType_COMPONENT_TYPE_TORSWITCH:
		return "torswitch"
	case rlav1.ComponentType_COMPONENT_TYPE_UMS:
		return "ums"
	case rlav1.ComponentType_COMPONENT_TYPE_CDU:
		return "cdu"
	default:
		return "unknown"
	}
}

// NewAPITray creates an APITray from the RLA protobuf Component
func NewAPITray(comp *rlav1.Component) *APITray {
	if comp == nil {
		return nil
	}

	apiTray := &APITray{
		Type:            componentTypeToString(comp.GetType()),
		FirmwareVersion: comp.GetFirmwareVersion(),
		ComponentID:     comp.GetComponentId(),
	}

	// Get info from DeviceInfo
	if comp.GetInfo() != nil {
		info := comp.GetInfo()
		if info.GetId() != nil {
			apiTray.ID = info.GetId().GetId()
		}
		apiTray.Name = info.GetName()
		apiTray.Manufacturer = info.GetManufacturer()
		if info.Model != nil {
			apiTray.Model = *info.Model
		}
		apiTray.SerialNumber = info.GetSerialNumber()
		if info.Description != nil {
			apiTray.Description = *info.Description
		}
	}

	// Get position
	if comp.GetPosition() != nil {
		pos := comp.GetPosition()
		apiTray.Position = &APITrayPosition{
			SlotID:  pos.GetSlotId(),
			TrayIdx: pos.GetTrayIdx(),
			HostID:  pos.GetHostId(),
		}
	}

	// Get rack ID
	if comp.GetRackId() != nil {
		apiTray.RackID = comp.GetRackId().GetId()
	}

	return apiTray
}

// NewAPITrays creates a slice of APITray from the RLA protobuf response
func NewAPITrays(resp *rlav1.GetComponentsResponse) []*APITray {
	if resp == nil {
		return []*APITray{}
	}

	trays := make([]*APITray, 0, len(resp.GetComponents()))
	for _, comp := range resp.GetComponents() {
		trays = append(trays, NewAPITray(comp))
	}

	return trays
}

// NewAPITraysWithFilter creates a slice of APITray with client-side filtering
func NewAPITraysWithFilter(resp *rlav1.GetComponentsResponse, filter *TrayFilterInput) []*APITray {
	if resp == nil {
		return []*APITray{}
	}

	trays := make([]*APITray, 0, len(resp.GetComponents()))
	for _, comp := range resp.GetComponents() {
		// Apply client-side filters
		if !matchesTrayFilter(comp, filter) {
			continue
		}
		trays = append(trays, NewAPITray(comp))
	}

	return trays
}

// matchesTrayFilter checks if a component matches the filter criteria
func matchesTrayFilter(comp *rlav1.Component, filter *TrayFilterInput) bool {
	if filter == nil {
		return true
	}

	// Filter by tray UUIDs (id=uuid1,uuid2,uuid3)
	if len(filter.IDs) > 0 {
		found := false
		compID := ""
		if comp.GetInfo() != nil && comp.GetInfo().GetId() != nil {
			compID = comp.GetInfo().GetId().GetId()
		}
		for _, id := range filter.IDs {
			if compID == id {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by rack ID
	if filter.RackID != nil {
		if comp.GetRackId() == nil || comp.GetRackId().GetId() != *filter.RackID {
			return false
		}
	}

	// Filter by tray type (compute|switch|powershelf|...)
	if filter.Type != nil {
		want := strings.ToLower(*filter.Type)
		got := componentTypeToString(comp.GetType())
		if got != want {
			return false
		}
	}

	// Filter by slot
	if filter.Slot != nil {
		if comp.GetPosition() == nil || comp.GetPosition().GetSlotId() != *filter.Slot {
			return false
		}
	}

	// Filter by index (tray_idx)
	if filter.Index != nil {
		if comp.GetPosition() == nil || comp.GetPosition().GetTrayIdx() != *filter.Index {
			return false
		}
	}

	// Filter by component IDs (componentId=id1,id2,id3)
	if len(filter.ComponentIDs) > 0 {
		found := false
		for _, cid := range filter.ComponentIDs {
			if comp.GetComponentId() == cid {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
