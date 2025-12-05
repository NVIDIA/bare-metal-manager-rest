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

// APIOperatingSystemSiteAssociation is the data structure to capture API representation of an sshkey association
type APIOperatingSystemSiteAssociation struct {
	// Site is the summary of the Site
	Site *APISiteSummary `json:"site"`
	// Version is the version of corresponding image on Site
	Version *string `json:"version"`
	// Status is the status of the OperatingSystemSiteAssociation
	Status string `json:"status"`
	// Created indicates the ISO datetime string for when the site was created
	Created time.Time `json:"created"`
	// Updated indicates the ISO datetime string for when the site was last updated
	Updated time.Time `json:"updated"`
}

// NewAPIOperatingSystemSiteAssociation accepts a DB layer OperatingSystemSiteAssociation object and returns an API object
func NewAPIOperatingSystemSiteAssociation(dbossa *cdbm.OperatingSystemSiteAssociation, ts *cdbm.TenantSite) *APIOperatingSystemSiteAssociation {
	apiossa := &APIOperatingSystemSiteAssociation{
		Version: dbossa.Version,
		Status:  dbossa.Status,
		Created: dbossa.Created,
		Updated: dbossa.Updated,
	}

	if dbossa.Site != nil {
		apiossa.Site = NewAPISiteSummary(dbossa.Site)
	}

	return apiossa
}
