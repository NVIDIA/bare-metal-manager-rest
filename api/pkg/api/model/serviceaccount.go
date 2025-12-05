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
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

// type APIServiceAccount is the data structure to capture API representation of a Service Account
type APIServiceAccount struct {
	// Enabled is a flag to indicate if the Service Account is enabled
	Enabled bool `json:"enabled"`
	// InfrastructureProviderID is the ID of the InfrastructureProvider
	InfrastructureProviderID *string `json:"infrastructureProviderId"`
	// ID is the unique UUID v4 identifier for the Service Account
	TenantID *string `json:"tenantId"`
}

// NewAPIServiceAccount accepts a DB layer ServiceAccount object and returns an API object
func NewAPIServiceAccount(serviceAccountEnabled bool, dbProvider *cdbm.InfrastructureProvider, dbTenant *cdbm.Tenant) *APIServiceAccount {
	apiServiceAccount := APIServiceAccount{
		Enabled: serviceAccountEnabled,
	}

	if dbProvider != nil {
		apiServiceAccount.InfrastructureProviderID = cdb.GetStrPtr(dbProvider.ID.String())
	}
	if dbTenant != nil {
		apiServiceAccount.TenantID = cdb.GetStrPtr(dbTenant.ID.String())
	}

	return &apiServiceAccount
}
