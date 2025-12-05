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
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

func TestNewAPIMachineInstanceType(t *testing.T) {
	type args struct {
		dbmit *cdbm.MachineInstanceType
	}

	dbmit := &cdbm.MachineInstanceType{
		ID:             uuid.New(),
		MachineID:      uuid.NewString(),
		InstanceTypeID: uuid.New(),
		Created:        time.Now(),
		Updated:        time.Now(),
	}

	tests := []struct {
		name string
		args args
		want *APIMachineInstanceType
	}{
		{
			name: "test new API Machine Instance Type initializer",
			args: args{
				dbmit: dbmit,
			},
			want: &APIMachineInstanceType{
				ID:             dbmit.ID.String(),
				MachineID:      dbmit.MachineID,
				InstanceTypeID: dbmit.InstanceTypeID.String(),
				Created:        dbmit.Created,
				Updated:        dbmit.Updated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPIMachineInstanceType(tt.args.dbmit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIMachineInstanceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIMachineInstanceTypeCreateRequest_Validate(t *testing.T) {
	type fields struct {
		MachineIDs []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test valid Machine Instance Type request",
			fields: fields{
				MachineIDs: []string{"test-machine-id", uuid.NewString()},
			},
			wantErr: false,
		},
		{
			name: "test invalid Machine Instance Type request, empty MachineIDs",
			fields: fields{
				MachineIDs: []string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mitcr := APIMachineInstanceTypeCreateRequest{
				MachineIDs: tt.fields.MachineIDs,
			}
			if err := mitcr.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("APIMachineInstanceTypeCreateRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
