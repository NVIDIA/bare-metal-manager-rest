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

	cdb "github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/db"
	cdbm "github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/db/model"
	"github.com/stretchr/testify/assert"
)

func TestMachineCapability_NewAPIMachineCapability(t *testing.T) {
	dbmc := &cdbm.MachineCapability{
		Type:      cdbm.MachineCapabilityTypeCPU,
		Name:      "AMD Opteron Series x10",
		Frequency: cdb.GetStrPtr("3.0GHz"),
		Capacity:  cdb.GetStrPtr("3.0GHz"),
		Vendor:    cdb.GetStrPtr("AMD"),
		Count:     cdb.GetIntPtr(2),
	}

	apimc := NewAPIMachineCapability(dbmc)
	assert.Equal(t, dbmc.Type, apimc.Type)
	assert.Equal(t, dbmc.Name, apimc.Name)
	assert.Equal(t, *dbmc.Frequency, *apimc.Frequency)
	assert.Equal(t, *dbmc.Capacity, *apimc.Capacity)
	assert.Equal(t, *dbmc.Vendor, *apimc.Vendor)
	assert.Equal(t, *dbmc.Count, *apimc.Count)
}
