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
	"testing"

	rlav1 "github.com/nvidia/carbide-rest/workflow-schema/rla/protobuf/v1"
	"github.com/stretchr/testify/assert"
)

// Helper function to create string pointer
func strPtr(s string) *string {
	return &s
}

// Helper function to create int32 pointer
func int32Ptr(i int32) *int32 {
	return &i
}

func TestMatchesTrayFilter(t *testing.T) {
	tests := []struct {
		name   string
		comp   *rlav1.Component
		filter *TrayFilterInput
		want   bool
	}{
		{
			name:   "nil filter matches all",
			comp:   &rlav1.Component{},
			filter: nil,
			want:   true,
		},
		{
			name:   "empty filter matches all",
			comp:   &rlav1.Component{},
			filter: &TrayFilterInput{},
			want:   true,
		},
		{
			name: "slot filter matches",
			comp: &rlav1.Component{
				Position: &rlav1.RackPosition{SlotId: 5},
			},
			filter: &TrayFilterInput{Slot: int32Ptr(5)},
			want:   true,
		},
		{
			name: "slot filter does not match",
			comp: &rlav1.Component{
				Position: &rlav1.RackPosition{SlotId: 3},
			},
			filter: &TrayFilterInput{Slot: int32Ptr(5)},
			want:   false,
		},
		{
			name:   "slot filter - component has no position",
			comp:   &rlav1.Component{},
			filter: &TrayFilterInput{Slot: int32Ptr(5)},
			want:   false,
		},
		{
			name: "index filter matches",
			comp: &rlav1.Component{
				Position: &rlav1.RackPosition{TrayIdx: 2},
			},
			filter: &TrayFilterInput{Index: int32Ptr(2)},
			want:   true,
		},
		{
			name: "index filter does not match",
			comp: &rlav1.Component{
				Position: &rlav1.RackPosition{TrayIdx: 1},
			},
			filter: &TrayFilterInput{Index: int32Ptr(2)},
			want:   false,
		},
		{
			name: "componentId filter matches",
			comp: &rlav1.Component{
				ComponentId: "comp-123",
			},
			filter: &TrayFilterInput{ComponentIDs: []string{"comp-123", "comp-456"}},
			want:   true,
		},
		{
			name: "componentId filter does not match",
			comp: &rlav1.Component{
				ComponentId: "comp-789",
			},
			filter: &TrayFilterInput{ComponentIDs: []string{"comp-123", "comp-456"}},
			want:   false,
		},
		{
			name: "combined filters - all match",
			comp: &rlav1.Component{
				ComponentId: "comp-123",
				Position: &rlav1.RackPosition{
					SlotId:  5,
					TrayIdx: 2,
				},
			},
			filter: &TrayFilterInput{
				Slot:         int32Ptr(5),
				Index:        int32Ptr(2),
				ComponentIDs: []string{"comp-123"},
			},
			want: true,
		},
		{
			name: "combined filters - one does not match",
			comp: &rlav1.Component{
				ComponentId: "comp-123",
				Position: &rlav1.RackPosition{
					SlotId:  5,
					TrayIdx: 1,
				},
			},
			filter: &TrayFilterInput{
				Slot:         int32Ptr(5),
				Index:        int32Ptr(2),
				ComponentIDs: []string{"comp-123"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesTrayFilter(tt.comp, tt.filter)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewAPITraysWithFilter(t *testing.T) {
	tests := []struct {
		name          string
		resp          *rlav1.GetComponentsResponse
		filter        *TrayFilterInput
		wantTrayCount int
	}{
		{
			name:          "nil response returns empty list",
			resp:          nil,
			filter:        nil,
			wantTrayCount: 0,
		},
		{
			name: "nil filter returns all trays",
			resp: &rlav1.GetComponentsResponse{
				Components: []*rlav1.Component{
					{
						Type: rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
						Info: &rlav1.DeviceInfo{Id: &rlav1.UUID{Id: "tray-1"}},
					},
					{
						Type: rlav1.ComponentType_COMPONENT_TYPE_NVLSWITCH,
						Info: &rlav1.DeviceInfo{Id: &rlav1.UUID{Id: "tray-2"}},
					},
				},
				Total: 2,
			},
			filter:        nil,
			wantTrayCount: 2,
		},
		{
			name: "filter by slot",
			resp: &rlav1.GetComponentsResponse{
				Components: []*rlav1.Component{
					{
						Type:     rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
						Info:     &rlav1.DeviceInfo{Id: &rlav1.UUID{Id: "tray-1"}},
						Position: &rlav1.RackPosition{SlotId: 1},
					},
					{
						Type:     rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
						Info:     &rlav1.DeviceInfo{Id: &rlav1.UUID{Id: "tray-2"}},
						Position: &rlav1.RackPosition{SlotId: 2},
					},
					{
						Type:     rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
						Info:     &rlav1.DeviceInfo{Id: &rlav1.UUID{Id: "tray-3"}},
						Position: &rlav1.RackPosition{SlotId: 1},
					},
				},
				Total: 3,
			},
			filter:        &TrayFilterInput{Slot: int32Ptr(1)},
			wantTrayCount: 2,
		},
		{
			name: "filter by componentId",
			resp: &rlav1.GetComponentsResponse{
				Components: []*rlav1.Component{
					{
						Type:        rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
						Info:        &rlav1.DeviceInfo{Id: &rlav1.UUID{Id: "tray-1"}},
						ComponentId: "comp-a",
					},
					{
						Type:        rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
						Info:        &rlav1.DeviceInfo{Id: &rlav1.UUID{Id: "tray-2"}},
						ComponentId: "comp-b",
					},
					{
						Type:        rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
						Info:        &rlav1.DeviceInfo{Id: &rlav1.UUID{Id: "tray-3"}},
						ComponentId: "comp-c",
					},
				},
				Total: 3,
			},
			filter:        &TrayFilterInput{ComponentIDs: []string{"comp-a", "comp-c"}},
			wantTrayCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAPITraysWithFilter(tt.resp, tt.filter)
			assert.NotNil(t, got)
			assert.Equal(t, tt.wantTrayCount, len(got))
		})
	}
}

func TestComponentTypeToString(t *testing.T) {
	tests := []struct {
		name string
		ct   rlav1.ComponentType
		want string
	}{
		{
			name: "compute type",
			ct:   rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
			want: "compute",
		},
		{
			name: "nvlswitch type",
			ct:   rlav1.ComponentType_COMPONENT_TYPE_NVLSWITCH,
			want: "switch",
		},
		{
			name: "powershelf type",
			ct:   rlav1.ComponentType_COMPONENT_TYPE_POWERSHELF,
			want: "powershelf",
		},
		{
			name: "torswitch type",
			ct:   rlav1.ComponentType_COMPONENT_TYPE_TORSWITCH,
			want: "torswitch",
		},
		{
			name: "ums type",
			ct:   rlav1.ComponentType_COMPONENT_TYPE_UMS,
			want: "ums",
		},
		{
			name: "cdu type",
			ct:   rlav1.ComponentType_COMPONENT_TYPE_CDU,
			want: "cdu",
		},
		{
			name: "unknown type",
			ct:   rlav1.ComponentType_COMPONENT_TYPE_UNKNOWN,
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := componentTypeToString(tt.ct)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewAPITray(t *testing.T) {
	description := "Test tray description"
	model := "GB200"

	tests := []struct {
		name string
		comp *rlav1.Component
		want *APITray
	}{
		{
			name: "nil component returns nil",
			comp: nil,
			want: nil,
		},
		{
			name: "basic compute tray",
			comp: &rlav1.Component{
				Type: rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
				Info: &rlav1.DeviceInfo{
					Id:           &rlav1.UUID{Id: "tray-id-123"},
					Name:         "compute-tray-1",
					Manufacturer: "NVIDIA",
					Model:        &model,
					SerialNumber: "TSN001",
					Description:  &description,
				},
				FirmwareVersion: "2.1.0",
				ComponentId:     "carbide-machine-456",
				Position: &rlav1.RackPosition{
					SlotId:  1,
					TrayIdx: 0,
					HostId:  1,
				},
				RackId: &rlav1.UUID{Id: "rack-id-789"},
			},
			want: &APITray{
				ID:              "tray-id-123",
				ComponentID:     "carbide-machine-456",
				Type:            "compute",
				Name:            "compute-tray-1",
				Manufacturer:    "NVIDIA",
				Model:           "GB200",
				SerialNumber:    "TSN001",
				Description:     "Test tray description",
				FirmwareVersion: "2.1.0",
				Position: &APITrayPosition{
					SlotID:  1,
					TrayIdx: 0,
					HostID:  1,
				},
				RackID: "rack-id-789",
			},
		},
		{
			name: "switch tray without optional fields",
			comp: &rlav1.Component{
				Type: rlav1.ComponentType_COMPONENT_TYPE_NVLSWITCH,
				Info: &rlav1.DeviceInfo{
					Id:           &rlav1.UUID{Id: "switch-tray-id"},
					Name:         "switch-tray-1",
					Manufacturer: "NVIDIA",
					SerialNumber: "SSN001",
				},
				FirmwareVersion: "1.5.0",
				Position: &rlav1.RackPosition{
					SlotId:  24,
					TrayIdx: 1,
				},
			},
			want: &APITray{
				ID:              "switch-tray-id",
				Type:            "switch",
				Name:            "switch-tray-1",
				Manufacturer:    "NVIDIA",
				SerialNumber:    "SSN001",
				FirmwareVersion: "1.5.0",
				Position: &APITrayPosition{
					SlotID:  24,
					TrayIdx: 1,
					HostID:  0,
				},
			},
		},
		{
			name: "powershelf tray",
			comp: &rlav1.Component{
				Type: rlav1.ComponentType_COMPONENT_TYPE_POWERSHELF,
				Info: &rlav1.DeviceInfo{
					Id:           &rlav1.UUID{Id: "power-tray-id"},
					Name:         "powershelf-1",
					Manufacturer: "NVIDIA",
					SerialNumber: "PSN001",
				},
				Position: &rlav1.RackPosition{
					SlotId: 48,
				},
				RackId: &rlav1.UUID{Id: "rack-abc"},
			},
			want: &APITray{
				ID:           "power-tray-id",
				Type:         "powershelf",
				Name:         "powershelf-1",
				Manufacturer: "NVIDIA",
				SerialNumber: "PSN001",
				Position: &APITrayPosition{
					SlotID:  48,
					TrayIdx: 0,
					HostID:  0,
				},
				RackID: "rack-abc",
			},
		},
		{
			name: "tray with minimal info",
			comp: &rlav1.Component{
				Type: rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
				Info: &rlav1.DeviceInfo{
					Id: &rlav1.UUID{Id: "minimal-tray"},
				},
			},
			want: &APITray{
				ID:   "minimal-tray",
				Type: "compute",
			},
		},
		{
			name: "tray without info",
			comp: &rlav1.Component{
				Type:        rlav1.ComponentType_COMPONENT_TYPE_UMS,
				ComponentId: "ums-component-123",
			},
			want: &APITray{
				Type:        "ums",
				ComponentID: "ums-component-123",
			},
		},
		{
			name: "tray without position",
			comp: &rlav1.Component{
				Type: rlav1.ComponentType_COMPONENT_TYPE_CDU,
				Info: &rlav1.DeviceInfo{
					Id:   &rlav1.UUID{Id: "cdu-tray-id"},
					Name: "cdu-1",
				},
			},
			want: &APITray{
				ID:       "cdu-tray-id",
				Type:     "cdu",
				Name:     "cdu-1",
				Position: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAPITray(tt.comp)

			if tt.want == nil {
				assert.Nil(t, got)
				return
			}

			assert.NotNil(t, got)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.ComponentID, got.ComponentID)
			assert.Equal(t, tt.want.Type, got.Type)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Manufacturer, got.Manufacturer)
			assert.Equal(t, tt.want.Model, got.Model)
			assert.Equal(t, tt.want.SerialNumber, got.SerialNumber)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.Equal(t, tt.want.FirmwareVersion, got.FirmwareVersion)
			assert.Equal(t, tt.want.RackID, got.RackID)

			if tt.want.Position != nil {
				assert.NotNil(t, got.Position)
				assert.Equal(t, tt.want.Position.SlotID, got.Position.SlotID)
				assert.Equal(t, tt.want.Position.TrayIdx, got.Position.TrayIdx)
				assert.Equal(t, tt.want.Position.HostID, got.Position.HostID)
			} else {
				assert.Nil(t, got.Position)
			}
		})
	}
}

func TestNewAPITrays(t *testing.T) {
	tests := []struct {
		name          string
		resp          *rlav1.GetComponentsResponse
		wantTrayCount int
	}{
		{
			name:          "nil response returns empty list",
			resp:          nil,
			wantTrayCount: 0,
		},
		{
			name: "response with multiple trays",
			resp: &rlav1.GetComponentsResponse{
				Components: []*rlav1.Component{
					{
						Type: rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
						Info: &rlav1.DeviceInfo{
							Id:   &rlav1.UUID{Id: "tray-1"},
							Name: "Tray 1",
						},
					},
					{
						Type: rlav1.ComponentType_COMPONENT_TYPE_NVLSWITCH,
						Info: &rlav1.DeviceInfo{
							Id:   &rlav1.UUID{Id: "tray-2"},
							Name: "Tray 2",
						},
					},
					{
						Type: rlav1.ComponentType_COMPONENT_TYPE_POWERSHELF,
						Info: &rlav1.DeviceInfo{
							Id:   &rlav1.UUID{Id: "tray-3"},
							Name: "Tray 3",
						},
					},
				},
				Total: 3,
			},
			wantTrayCount: 3,
		},
		{
			name: "empty components list",
			resp: &rlav1.GetComponentsResponse{
				Components: []*rlav1.Component{},
				Total:      0,
			},
			wantTrayCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAPITrays(tt.resp)

			assert.NotNil(t, got)
			assert.Equal(t, tt.wantTrayCount, len(got))

			// Verify each tray was converted correctly
			if tt.resp != nil {
				for i, comp := range tt.resp.Components {
					if comp.Info != nil && comp.Info.Id != nil {
						assert.Equal(t, comp.Info.Id.Id, got[i].ID)
					}
					if comp.Info != nil {
						assert.Equal(t, comp.Info.Name, got[i].Name)
					}
				}
			}
		})
	}
}
