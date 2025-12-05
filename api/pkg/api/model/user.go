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

// APIUser is a data structure to capture information about user at the API layer
type APIUser struct {
	// ID is the unique UUID v4 identifier of the user in Forge Cloud
	ID string `json:"id"`
	// FirstName denotes the first name of the user
	FirstName *string `json:"firstName"`
	// LastName denotes the surname of the user
	LastName *string `json:"lastName"`
	// Email contains the email used by the user to register with NGC
	Email *string `json:"email"`
	// Created indicates the ISO datetime string for when the user was created in Forge
	Created time.Time `json:"created"`
	// Updated indicates the ISO datetime string for when the user was last updated in Forge
	Updated time.Time `json:"updated"`
}

// NewAPIUserFromDBUser creates and returns a new APIUser object
func NewAPIUserFromDBUser(dbUser cdbm.User) *APIUser {
	apiUser := &APIUser{
		ID:        dbUser.ID.String(),
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		Email:     dbUser.Email,
		Created:   dbUser.Created,
		Updated:   dbUser.Updated,
	}

	return apiUser
}
