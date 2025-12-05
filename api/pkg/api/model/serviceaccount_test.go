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
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

func TestNewAPIServiceAccount(t *testing.T) {
	type args struct {
		serviceAccountEnabled bool
		dbProvider            *cdbm.InfrastructureProvider
		dbTenant              *cdbm.Tenant
	}

	dbProvider := &cdbm.InfrastructureProvider{
		ID: uuid.New(),
	}
	dbTenant := &cdbm.Tenant{
		ID: uuid.New(),
	}

	tests := []struct {
		name string
		args args
		want *APIServiceAccount
	}{
		{
			name: "test NewAPIServiceAccount with service account enabled",
			args: args{
				serviceAccountEnabled: true,
				dbProvider:            dbProvider,
				dbTenant:              dbTenant,
			},
			want: &APIServiceAccount{
				Enabled:                  true,
				InfrastructureProviderID: cdb.GetStrPtr(dbProvider.ID.String()),
				TenantID:                 cdb.GetStrPtr(dbTenant.ID.String()),
			},
		},
		{
			name: "test NewAPIServiceAccount with service account disabled",
			args: args{
				serviceAccountEnabled: false,
			},
			want: &APIServiceAccount{
				Enabled: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAPIServiceAccount(tt.args.serviceAccountEnabled, tt.args.dbProvider, tt.args.dbTenant)
			assert.Equal(t, tt.want, got)
		})
	}
}
