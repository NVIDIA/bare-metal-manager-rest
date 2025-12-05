// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package model

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

func TestNewAPINVLinkInterface(t *testing.T) {
	type args struct {
		dbnvli *cdbm.NVLinkInterface
	}

	dbnvli := &cdbm.NVLinkInterface{
		ID:                       uuid.New(),
		InstanceID:               uuid.New(),
		SiteID:                   uuid.New(),
		NVLinkLogicalPartitionID: uuid.New(),
		DeviceInstance:           1,
		Status:                   cdbm.NVLinkInterfaceStatusReady,
		Created:                  time.Now(),
		Updated:                  time.Now(),
	}

	tests := []struct {
		name string
		args args
		want *APINVLinkInterface
	}{
		{
			name: "test new API NVLink Interface initializer",
			args: args{
				dbnvli: dbnvli,
			},
			want: &APINVLinkInterface{
				ID:                       dbnvli.ID.String(),
				InstanceID:               dbnvli.InstanceID.String(),
				NVLinkLogicalPartitionID: dbnvli.NVLinkLogicalPartitionID.String(),
				DeviceInstance:           dbnvli.DeviceInstance,
				Status:                   dbnvli.Status,
				Created:                  dbnvli.Created,
				Updated:                  dbnvli.Updated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPINVLinkInterface(tt.args.dbnvli); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPINVLinkInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPINVLinkInterfaceCreateOrUpdateRequest_Validate(t *testing.T) {
	type fields struct {
		nvLinkLogicalPartitionID string
		deviceInstance           int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test validation success",
			fields: fields{
				nvLinkLogicalPartitionID: uuid.New().String(),
				deviceInstance:           0,
			},
			wantErr: false,
		},
		{
			name: "test validation failure, invalid NVLink Logical Partition ID",
			fields: fields{
				nvLinkLogicalPartitionID: "badid",
				deviceInstance:           1,
			},
			wantErr: true,
		},
		{
			name: "test validation failure, GPU Index not supported",
			fields: fields{
				nvLinkLogicalPartitionID: uuid.New().String(),
				deviceInstance:           4,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nvlicr := APINVLinkInterfaceCreateOrUpdateRequest{
				NVLinkLogicalPartitionID: tt.fields.nvLinkLogicalPartitionID,
				DeviceInstance:           tt.fields.deviceInstance,
			}
			err := nvlicr.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
