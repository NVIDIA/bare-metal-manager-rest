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
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

func TestNewAPISSHKeyAssociation(t *testing.T) {
	ska := cdbm.SSHKeyAssociation{
		ID:            uuid.New(),
		SSHKeyID:      uuid.New(),
		SSHKeyGroupID: uuid.New(),
		Created:       time.Now(),
		Updated:       time.Now(),
	}
	apiska := NewAPISSHKeyAssociation(&ska)
	assert.Equal(t, apiska.ID, ska.ID.String())
	assert.Equal(t, apiska.SSHKeyID, ska.SSHKeyID.String())
	assert.Equal(t, apiska.SSHKeyGroupID, ska.SSHKeyGroupID.String())
}
