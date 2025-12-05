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

func TestNewAPIInfrastructureProvider(t *testing.T) {
	type args struct {
		dbip *cdbm.InfrastructureProvider
	}

	dbip := &cdbm.InfrastructureProvider{
		ID:             uuid.New(),
		Name:           "test-infrastructure-provider",
		DisplayName:    nil,
		Org:            "test-org",
		OrgDisplayName: cdb.GetStrPtr("Org Display name"),
		Created:        time.Now(),
		Updated:        time.Now(),
	}

	ipAPIInfrastructureProvider := APIInfrastructureProvider{
		ID:             dbip.ID.String(),
		Org:            dbip.Org,
		OrgDisplayName: dbip.OrgDisplayName,
		Created:        dbip.Created,
		Updated:        dbip.Updated,
	}

	tests := []struct {
		name string
		args args
		want *APIInfrastructureProvider
	}{
		{
			name: "test initializing API model for Infrastructure Provider",
			args: args{
				dbip: dbip,
			},
			want: &ipAPIInfrastructureProvider,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPIInfrastructureProvider(tt.args.dbip); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIInfrastructureProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAPIInfrastructureProviderSummary(t *testing.T) {
	dbip := &cdbm.InfrastructureProvider{
		ID:             uuid.New(),
		Name:           "test-infrastructure-provider",
		DisplayName:    nil,
		Org:            "test-org",
		OrgDisplayName: cdb.GetStrPtr("Org Display name"),
		Created:        time.Now(),
		Updated:        time.Now(),
	}

	type args struct {
		dbip *cdbm.InfrastructureProvider
	}
	tests := []struct {
		name string
		args args
		want *APIInfrastructureProviderSummary
	}{
		{
			name: "test init API summary model for Infrastructure Provider",
			args: args{
				dbip: dbip,
			},
			want: &APIInfrastructureProviderSummary{
				Org:            dbip.Org,
				OrgDisplayName: dbip.OrgDisplayName,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPIInfrastructureProviderSummary(tt.args.dbip); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIInfrastructureProviderSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}
