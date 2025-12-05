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
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

// APIMachineInstanceTypeCreateRequest is the data structure to capture user request to create a new MachineInstanceType
type APIMachineInstanceTypeCreateRequest struct {
	// MachineID is the ID of the Machine
	MachineIDs []string `json:"machineIds"`
}

// Validate ensure the values passed in request are acceptable
func (mitcr APIMachineInstanceTypeCreateRequest) Validate() error {
	err := validation.ValidateStruct(&mitcr,
		validation.Field(&mitcr.MachineIDs,
			validation.Required.Error("at least one machine ID is required"),
		),
	)
	return err
}

// APIMachineInstanceType is the data structure to capture Machine Instance Type
type APIMachineInstanceType struct {
	// ID is the unique UUID v4 identifier for the Machine Instance Type
	ID string `json:"id"`
	// MachineID is the ID of the associated Machine
	MachineID string `json:"machineId"`
	// InstanceTypeID is the ID of the associated Instance Type
	InstanceTypeID string `json:"instanceTypeId"`
	// Created is the date and time the Machine Instance Type was created
	Created time.Time `json:"created"`
	// Updated is the date and time the Machine Instance Type was last updated
	Updated time.Time `json:"updated"`
}

// NewAPIMachineInstanceType creates a new APIMachineInstanceType
func NewAPIMachineInstanceType(dbmit *cdbm.MachineInstanceType) *APIMachineInstanceType {
	return &APIMachineInstanceType{
		ID:             dbmit.ID.String(),
		MachineID:      dbmit.MachineID,
		InstanceTypeID: dbmit.InstanceTypeID.String(),
		Created:        dbmit.Created,
		Updated:        dbmit.Updated,
	}
}
