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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

func TestMachine_NewAPIFabric(t *testing.T) {
	fguid := "Ifabric01"
	dbf := &cdbm.Fabric{
		ID:                       fguid,
		InfrastructureProviderID: uuid.New(),
		SiteID:                   uuid.New(),
		Status:                   cdbm.FabricStatusPending,
		Created:                  cdb.GetCurTime(),
		Updated:                  cdb.GetCurTime(),
	}

	dbsds := []cdbm.StatusDetail{
		{
			ID:       uuid.New(),
			EntityID: dbf.ID,
			Status:   dbf.Status,
			Created:  cdb.GetCurTime(),
			Updated:  cdb.GetCurTime(),
		},
	}

	dbf.Site = &cdbm.Site{
		ID:                       dbf.SiteID,
		Name:                     "test-site",
		Description:              cdb.GetStrPtr("Test Description"),
		InfrastructureProviderID: dbf.InfrastructureProviderID,
		Status:                   cdbm.SiteStatusRegistered,
		Created:                  cdb.GetCurTime(),
		Updated:                  cdb.GetCurTime(),
		CreatedBy:                uuid.New(),
	}

	apif := NewAPIFabric(dbf, dbsds)
	assert.NotNil(t, apif)

	assert.Equal(t, apif.ID, dbf.ID)
	assert.Equal(t, apif.InfrastructureProviderID, dbf.InfrastructureProviderID.String())
	assert.Equal(t, apif.SiteID, dbf.SiteID.String())
	assert.NotNil(t, apif.Site)
	assert.Equal(t, apif.Site.Name, dbf.Site.Name)

	for i, v := range dbsds {
		assert.Equal(t, apif.StatusHistory[i].Status, v.Status)
	}
}
