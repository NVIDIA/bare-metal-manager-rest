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
			name: "returns all trays",
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
			wantTrayCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAPITrays(tt.resp)
			assert.NotNil(t, got)
			assert.Equal(t, tt.wantTrayCount, len(got))
		})
	}
}

func TestProtoToAPIComponentTypeName(t *testing.T) {
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
			got := ProtoToAPIComponentTypeName[rlav1.ComponentType_name[int32(tt.ct)]]
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

func TestAPITrayPosition_FromProto(t *testing.T) {
	pos := &APITrayPosition{}
	pos.FromProto(&rlav1.RackPosition{SlotId: 2, TrayIdx: 1, HostId: 0})
	assert.Equal(t, int32(2), pos.SlotID)
	assert.Equal(t, int32(1), pos.TrayIdx)
	assert.Equal(t, int32(0), pos.HostID)

	pos.FromProto(nil) // no-op
	assert.Equal(t, int32(2), pos.SlotID)
}

func TestAPITray_FromProto(t *testing.T) {
	comp := &rlav1.Component{
		Type:            rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
		ComponentId:     "comp-1",
		FirmwareVersion: "1.0",
		Info: &rlav1.DeviceInfo{
			Id:   &rlav1.UUID{Id: "tray-uuid"},
			Name: "My Tray",
		},
		Position: &rlav1.RackPosition{SlotId: 3, TrayIdx: 0, HostId: 1},
		RackId:   &rlav1.UUID{Id: "rack-uuid"},
	}
	at := &APITray{}
	at.FromProto(comp)
	assert.Equal(t, "compute", at.Type)
	assert.Equal(t, "comp-1", at.ComponentID)
	assert.Equal(t, "tray-uuid", at.ID)
	assert.Equal(t, "My Tray", at.Name)
	assert.Equal(t, "rack-uuid", at.RackID)
	assert.NotNil(t, at.Position)
	assert.Equal(t, int32(3), at.Position.SlotID)
	assert.Equal(t, int32(0), at.Position.TrayIdx)
	assert.Equal(t, int32(1), at.Position.HostID)

	at.FromProto(nil) // no-op, fields unchanged
	assert.Equal(t, "tray-uuid", at.ID)
}

func TestGetProtoTrayOrderByFromQueryParam(t *testing.T) {
	tests := []struct {
		field     string
		direction string
		wantNil   bool
	}{
		{"name", "ASC", false},
		{"manufacturer", "DESC", false},
		{"model", "ASC", false},
		{"type", "DESC", false},
		{"invalid", "ASC", true},
		{"name", "asc", false},
	}
	for _, tt := range tests {
		t.Run(tt.field+"_"+tt.direction, func(t *testing.T) {
			got := GetProtoTrayOrderByFromQueryParam(tt.field, tt.direction)
			if tt.wantNil {
				assert.Nil(t, got)
				return
			}
			assert.NotNil(t, got)
			assert.Equal(t, tt.direction, got.Direction)
			assert.NotNil(t, got.GetComponentField())
		})
	}
}
