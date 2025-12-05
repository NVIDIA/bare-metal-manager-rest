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
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	validationis "github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/nvidia/carbide-rest/api/pkg/api/model/util"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

var (
	// ErrValidationInfiniBandPartitionAssociation is the error when no associations are specified in the security group
	ErrValidationInfiniBandPartitionAssociation = errors.New("at least one security group association is required")
)

// APIInfiniBandPartitionCreateRequest is the data structure to capture instance request to create a new InfiniBandPartition
type APIInfiniBandPartitionCreateRequest struct {
	// Name is the name of the InfiniBand Partition
	Name string `json:"name"`
	// Description is the description of the InfiniBand Partition
	Description *string `json:"description"`
	// SiteID is the ID of the Site
	SiteID string `json:"siteId"`
}

// Validate ensure the values passed in request are acceptable
func (fbcr APIInfiniBandPartitionCreateRequest) Validate() error {
	err := validation.ValidateStruct(&fbcr,
		validation.Field(&fbcr.Name,
			validation.Required.Error(validationErrorStringLength),
			validation.By(util.ValidateNameCharacters),
			validation.Length(2, 256).Error(validationErrorStringLength)),
		validation.Field(&fbcr.SiteID,
			validation.Required.Error(validationErrorValueRequired),
			validationis.UUID.Error(validationErrorInvalidUUID)),
	)
	if err != nil {
		return err
	}
	return nil
}

// APIInfiniBandPartitionUpdateRequest is the data structure to capture user request to update a InfiniBandPartition
type APIInfiniBandPartitionUpdateRequest struct {
	// Name is the name of the InfiniBand Partition
	Name *string `json:"name"`
	// Description is the description of the InfiniBand Partition
	Description *string `json:"description"`
}

// Validate ensure the values passed in request are acceptable
func (fbur APIInfiniBandPartitionUpdateRequest) Validate() error {
	return validation.ValidateStruct(&fbur,
		validation.Field(&fbur.Name,
			validation.When(fbur.Name != nil, validation.Required.Error(validationErrorStringLength)),
			validation.When(fbur.Name != nil, validation.By(util.ValidateNameCharacters)),
			validation.When(fbur.Name != nil, validation.Length(2, 256).Error(validationErrorStringLength))),
	)
}

// APIInfiniBandPartition is the data structure to capture API representation of a InfiniBand Partition
type APIInfiniBandPartition struct {
	// ID is the unique UUID v4 identifier for the InfiniBand Partition
	ID string `json:"id"`
	// Name is the name of the InfiniBand Partition
	Name string `json:"name"`
	// Description is the description of the InfiniBand Partition
	Description *string `json:"description"`
	// SiteID is the ID of the Site
	SiteID string `json:"siteId"`
	// Site is the summary of the Site
	Site *APISiteSummary `json:"site,omitempty"`
	// TenantID is the ID of the Tenant
	TenantID string `json:"tenantId"`
	// Tenant is the summary of the tenant
	Tenant *APITenantSummary `json:"tenant,omitempty"`
	// Controller IB Partition ID is the ID of the Site Controller IB partition
	ControllerIBPartitionID *string `json:"controllerIBPartitionId"`
	// Partition Key is the key of IB partition
	PartitionKey *string `json:"partitionKey"`
	// Partition Name is the name of IB partition
	PartitionName *string `json:"partitionName"`
	// Service Level is the service level of IB partition
	ServiceLevel *int `json:"serviceLevel"`
	// Rate Limit is the rate limit of IB partition
	RateLimit *float32 `json:"rateLimit"`
	// Mtu of the IB partition
	Mtu *int `json:"mtu"`
	// EnableSharp indicates if sharp enable on the IB partition or not
	EnableSharp *bool `json:"enableSharp"`
	// Status is the status o the InfiniBand Partition
	Status string `json:"status"`
	// StatusHistory is the status detail records for the InfiniBand Partition over time
	StatusHistory []APIStatusDetail `json:"statusHistory"`
	// Created indicates the ISO datetime string for when the InfiniBand Partition was created
	Created time.Time `json:"created"`
	// Updated indicates the ISO datetime string for when the InfiniBand Partition was last updated
	Updated time.Time `json:"updated"`
}

// NewAPIInfiniBandPartition accepts a DB layer InfiniBandPartition object and returns an API object
func NewAPIInfiniBandPartition(dibp *cdbm.InfiniBandPartition, dbsds []cdbm.StatusDetail) *APIInfiniBandPartition {
	apiibp := &APIInfiniBandPartition{
		ID:            dibp.ID.String(),
		Name:          dibp.Name,
		Description:   dibp.Description,
		SiteID:        dibp.SiteID.String(),
		TenantID:      dibp.TenantID.String(),
		PartitionKey:  dibp.PartitionKey,
		PartitionName: dibp.PartitionName,
		ServiceLevel:  dibp.ServiceLevel,
		RateLimit:     dibp.RateLimit,
		Mtu:           dibp.Mtu,
		EnableSharp:   dibp.EnableSharp,
		Status:        dibp.Status,
		Created:       dibp.Created,
		Updated:       dibp.Updated,
	}

	if dibp.ControllerIBPartitionID != nil {
		apiibp.ControllerIBPartitionID = util.GetUUIDPtrToStrPtr(dibp.ControllerIBPartitionID)
	}

	if dibp.Site != nil {
		apiibp.Site = NewAPISiteSummary(dibp.Site)
	}

	if dibp.Tenant != nil {
		apiibp.Tenant = NewAPITenantSummary(dibp.Tenant)
	}

	apiibp.StatusHistory = []APIStatusDetail{}
	for _, dbsd := range dbsds {
		apiibp.StatusHistory = append(apiibp.StatusHistory, NewAPIStatusDetail(dbsd))
	}
	return apiibp
}

// APIInfiniBandPartitionSummary is the data structure to capture API summary of a InfiniBandPartition
type APIInfiniBandPartitionSummary struct {
	// ID of the InfiniBand Partition
	ID string `json:"id"`
	// Name of the InfiniBand Partition
	Name string `json:"name"`
	// SiteID is the ID of the Site
	SiteID string `json:"siteId"`
	// Controller IB Partition is the ID of the Site Controller Partition corresponding to the InfiniBand Partition
	ControllerIBPartitionID *string `json:"controllerIBPartitionId"`
	// Status is the status of the InfiniBand Partition
	Status string `json:"status"`
}

// NewAPIInfiniBandPartitionSummary accepts a DB layer InfiniBandPartition object returns an API layer object
func NewAPIInfiniBandPartitionSummary(dbibp *cdbm.InfiniBandPartition) *APIInfiniBandPartitionSummary {
	apiibps := APIInfiniBandPartitionSummary{
		ID:     dbibp.ID.String(),
		Name:   dbibp.Name,
		SiteID: dbibp.SiteID.String(),
		Status: dbibp.Status,
	}
	if dbibp.ControllerIBPartitionID != nil {
		apiibps.ControllerIBPartitionID = util.GetUUIDPtrToStrPtr(dbibp.ControllerIBPartitionID)
	}
	return &apiibps
}
