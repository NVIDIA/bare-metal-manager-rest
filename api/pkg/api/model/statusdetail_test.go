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
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

func TestNewAPIStatusDetail(t *testing.T) {
	type args struct {
		dbsd cdbm.StatusDetail
	}

	dbsd := cdbm.StatusDetail{
		ID:       uuid.New(),
		EntityID: uuid.NewString(),
		Status:   cdbm.SiteStatusPending,
		Message:  cdb.GetStrPtr("received request, pending processing"),
		Count:    1,
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	tests := []struct {
		name string
		args args
		want APIStatusDetail
	}{
		{
			name: "get new APIStatusDetail",
			args: args{
				dbsd: dbsd,
			},
			want: APIStatusDetail{
				Status:  dbsd.Status,
				Message: dbsd.Message,
				Created: dbsd.Created,
				Updated: dbsd.Updated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPIStatusDetail(tt.args.dbsd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIStatusDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
