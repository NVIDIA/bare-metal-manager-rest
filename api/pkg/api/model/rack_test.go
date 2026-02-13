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

	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewAPIRack(t *testing.T) {
	description := "Test rack description"
	model := "NVL72"

	tests := []struct {
		name           string
		rack           *rlav1.Rack
		withComponents bool
		want           *APIRack
	}{
		{
			name:           "nil rack returns nil",
			rack:           nil,
			withComponents: false,
			want:           nil,
		},
		{
			name: "basic rack without components",
			rack: &rlav1.Rack{
				Info: &rlav1.DeviceInfo{
					Id:           &rlav1.UUID{Id: "test-rack-id"},
					Name:         "test-rack",
					Manufacturer: "NVIDIA",
					Model:        &model,
					SerialNumber: "SN12345",
					Description:  &description,
				},
				Location: &rlav1.Location{
					Region:     "us-west-2",
					Datacenter: "DC1",
					Room:       "Room-A",
					Position:   "Row-1-Pos-5",
				},
			},
			withComponents: false,
			want: &APIRack{
				ID:           "test-rack-id",
				Name:         "test-rack",
				Manufacturer: "NVIDIA",
				Model:        "NVL72",
				SerialNumber: "SN12345",
				Description:  "Test rack description",
				Location: &APIRackLocation{
					Region:     "us-west-2",
					Datacenter: "DC1",
					Room:       "Room-A",
					Position:   "Row-1-Pos-5",
				},
				Components: nil,
			},
		},
		{
			name: "rack with components",
			rack: &rlav1.Rack{
				Info: &rlav1.DeviceInfo{
					Id:   &rlav1.UUID{Id: "rack-with-components"},
					Name: "rack-1",
				},
				Components: []*rlav1.Component{
					{
						Type: rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
						Info: &rlav1.DeviceInfo{
							Id:           &rlav1.UUID{Id: "comp-1"},
							Name:         "compute-node-1",
							SerialNumber: "CSN001",
							Manufacturer: "NVIDIA",
						},
						FirmwareVersion: "1.0.0",
						Position: &rlav1.RackPosition{
							SlotId: 1,
						},
						ComponentId: "carbide-machine-123",
					},
					{
						Type: rlav1.ComponentType_COMPONENT_TYPE_TORSWITCH,
						Info: &rlav1.DeviceInfo{
							Id:   &rlav1.UUID{Id: "comp-2"},
							Name: "switch-1",
						},
						Position: &rlav1.RackPosition{
							SlotId: 48,
						},
					},
				},
			},
			withComponents: true,
			want: &APIRack{
				ID:   "rack-with-components",
				Name: "rack-1",
				Components: []*APIRackComponent{
					{
						ID:              "comp-1",
						ComponentID:     "carbide-machine-123",
						Type:            "COMPONENT_TYPE_COMPUTE",
						Name:            "compute-node-1",
						SerialNumber:    "CSN001",
						Manufacturer:    "NVIDIA",
						FirmwareVersion: "1.0.0",
						SlotID:          1,
					},
					{
						ID:       "comp-2",
						Type:     "COMPONENT_TYPE_TORSWITCH",
						Name:     "switch-1",
						SlotID:   48,
					},
				},
			},
		},
		{
			name: "rack with components but withComponents=false",
			rack: &rlav1.Rack{
				Info: &rlav1.DeviceInfo{
					Id:   &rlav1.UUID{Id: "rack-id"},
					Name: "rack-name",
				},
				Components: []*rlav1.Component{
					{
						Type: rlav1.ComponentType_COMPONENT_TYPE_COMPUTE,
						Info: &rlav1.DeviceInfo{
							Id:   &rlav1.UUID{Id: "comp-1"},
							Name: "compute-node-1",
						},
					},
				},
			},
			withComponents: false,
			want: &APIRack{
				ID:         "rack-id",
				Name:       "rack-name",
				Components: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAPIRack(tt.rack, tt.withComponents)

			if tt.want == nil {
				assert.Nil(t, got)
				return
			}

			assert.NotNil(t, got)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Manufacturer, got.Manufacturer)
			assert.Equal(t, tt.want.Model, got.Model)
			assert.Equal(t, tt.want.SerialNumber, got.SerialNumber)
			assert.Equal(t, tt.want.Description, got.Description)

			if tt.want.Location != nil {
				assert.NotNil(t, got.Location)
				assert.Equal(t, tt.want.Location.Region, got.Location.Region)
				assert.Equal(t, tt.want.Location.Datacenter, got.Location.Datacenter)
				assert.Equal(t, tt.want.Location.Room, got.Location.Room)
				assert.Equal(t, tt.want.Location.Position, got.Location.Position)
			}

			if tt.want.Components != nil {
				assert.NotNil(t, got.Components)
				assert.Equal(t, len(tt.want.Components), len(got.Components))
				for i, wantComp := range tt.want.Components {
					gotComp := got.Components[i]
					assert.Equal(t, wantComp.ID, gotComp.ID)
					assert.Equal(t, wantComp.ComponentID, gotComp.ComponentID)
					assert.Equal(t, wantComp.Type, gotComp.Type)
					assert.Equal(t, wantComp.Name, gotComp.Name)
					assert.Equal(t, wantComp.SerialNumber, gotComp.SerialNumber)
					assert.Equal(t, wantComp.Manufacturer, gotComp.Manufacturer)
					assert.Equal(t, wantComp.FirmwareVersion, gotComp.FirmwareVersion)
					assert.Equal(t, wantComp.SlotID, gotComp.SlotID)
				}
			} else {
				assert.Nil(t, got.Components)
			}
		})
	}
}
