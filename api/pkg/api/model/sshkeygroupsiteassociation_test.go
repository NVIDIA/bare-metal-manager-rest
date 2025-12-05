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
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

func TestNewAPISSHKeyGroupSiteAssociation(t *testing.T) {
	skgsa := cdbm.SSHKeyGroupSiteAssociation{
		ID:            uuid.New(),
		SSHKeyGroupID: uuid.New(),
		SiteID:        uuid.New(),
		Version:       db.GetStrPtr("1234"),
		Status:        cdbm.SSHKeyGroupSiteAssociationStatusSyncing,
		Created:       time.Now(),
		Updated:       time.Now(),
	}
	apiskgsa := NewAPISSHKeyGroupSiteAssociation(&skgsa, nil)
	assert.Equal(t, apiskgsa.ControllerKeySetVersion, skgsa.Version)
	assert.Equal(t, apiskgsa.Status, skgsa.Status)
}
