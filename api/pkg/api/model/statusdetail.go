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

	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

// APIStatusDetail captures API representation of a status detail DB object
type APIStatusDetail struct {
	// Status denotes the state of the associated entity at a particular time
	Status string `json:"status"`
	// Message contains the description of the state and cause/remedy in case of error
	Message *string `json:"message"`
	// Created indicates the ISO datetime string for when the associated entity assumed the status
	Created time.Time `json:"created"`
	// Updated indicates the ISO datetime string for when the associated entity was last found to have this status
	Updated time.Time `json:"updated"`
}

// NewAPIStatusDetail creates an API status detail object from status detail DB entry
func NewAPIStatusDetail(dbsd cdbm.StatusDetail) APIStatusDetail {
	apiStatusDetail := APIStatusDetail{
		Status:  dbsd.Status,
		Message: dbsd.Message,
		Created: dbsd.Created,
		Updated: dbsd.Updated,
	}

	return apiStatusDetail
}
