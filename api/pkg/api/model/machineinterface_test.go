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
	"github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

func TestMachineInterface_NewAPIMachineInterface(t *testing.T) {
	dbmi := &cdbm.MachineInterface{
		ID:                    uuid.New(),
		MachineID:             uuid.NewString(),
		ControllerInterfaceID: db.GetUUIDPtr(uuid.New()),
		ControllerSegmentID:   db.GetUUIDPtr(uuid.New()),
		AttachedDPUMachineID:  db.GetStrPtr(uuid.NewString()),
		Hostname:              db.GetStrPtr("test.com"),
		IsPrimary:             true,
		SubnetID:              db.GetUUIDPtr(uuid.New()),
		MacAddress:            db.GetStrPtr("00:00:00:00:00:00"),
		IPAddresses:           []string{"192.168.0.1, 172.168.0.1"},
	}
	apimi := NewAPIMachineInterface(dbmi, true)
	assert.Equal(t, dbmi.ID.String(), apimi.ID)
	assert.Equal(t, dbmi.MachineID, apimi.MachineID)
	assert.Equal(t, dbmi.ControllerInterfaceID.String(), *apimi.ControllerInterfaceID)
	assert.Equal(t, dbmi.ControllerSegmentID.String(), *apimi.ControllerSegmentID)
	assert.Equal(t, *dbmi.AttachedDPUMachineID, *apimi.AttachedDPUMachineID)
	assert.Equal(t, dbmi.SubnetID.String(), *apimi.SubnetID)
	assert.Equal(t, *dbmi.Hostname, *apimi.Hostname)
	assert.Equal(t, dbmi.IsPrimary, apimi.IsPrimary)
	assert.Equal(t, *dbmi.MacAddress, *apimi.MacAddress)
	assert.Equal(t, len(dbmi.IPAddresses), len(apimi.IPAddresses))
	for i, v := range apimi.IPAddresses {
		assert.Equal(t, dbmi.IPAddresses[i], v)
	}
}
