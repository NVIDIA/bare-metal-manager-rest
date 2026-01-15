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

	"github.com/google/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	validationis "github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/nvidia/carbide-rest/api/pkg/api/model/util"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

const (
	// ExpectedMachineMaxBatchItems is the maximum number of ExpectedMachines allowed in a single batch operation
	ExpectedMachineMaxBatchItems = 100
)

// APIExpectedMachineCreateRequest is the data structure to capture instance request to create a new ExpectedMachine
type APIExpectedMachineCreateRequest struct {
	// SiteID is the ID of the Site
	SiteID string `json:"siteId"`
	// BmcMacAddress is the MAC address of the expected machine's BMC
	BmcMacAddress string `json:"bmcMacAddress"`
	// BmcUsername is the username of the expected machine's BMC
	DefaultBmcUsername *string `json:"defaultBmcUsername"`
	// DefaultBmcPassword is the password of the expected machine's BMC
	DefaultBmcPassword *string `json:"defaultBmcPassword"`
	// ChassisSerialNumber is the serial number of the expected machine's chassis
	ChassisSerialNumber string `json:"chassisSerialNumber"`
	// FallbackDPUSerialNumbers is the serial numbers of the expected machine's fallback DPUs
	FallbackDPUSerialNumbers []string `json:"fallbackDPUSerialNumbers"`
	// SkuId is the optional UUID for an SKU
	SkuID *string `json:"skuId"`
	// Labels is the labels of the expected machine
	Labels map[string]string `json:"labels"`
}

// Validate ensure the values passed in request are acceptable
func (emcr *APIExpectedMachineCreateRequest) Validate() error {
	err := validation.ValidateStruct(emcr,
		validation.Field(&emcr.BmcMacAddress,
			validation.Required.Error(validationErrorValueRequired),
			validationis.MAC),
		validation.Field(&emcr.DefaultBmcUsername,
			validation.Length(0, 16).Error("BMC username must be 16 characters or less")),
		validation.Field(&emcr.DefaultBmcPassword,
			validation.Length(0, 20).Error("BMC password must be 20 characters or less")),
		validation.Field(&emcr.ChassisSerialNumber,
			validation.Required.Error(validationErrorValueRequired),
			validation.Length(1, 32).Error("Chassis serial number must be 32 characters or less")),
		validation.Field(&emcr.SkuID,
			validation.NilOrNotEmpty.Error("SkuID cannot be empty")),
	)

	if err != nil {
		return err
	}

	// Labels validation
	if emcr.Labels != nil {
		if len(emcr.Labels) > util.LabelCountMax {
			return validation.Errors{
				"labels": util.ErrValidationLabelCount,
			}
		}

		for key, value := range emcr.Labels {
			if key == "" {
				return validation.Errors{
					"labels": util.ErrValidationLabelKeyEmpty,
				}
			}

			// Key validation
			if len(key) > util.LabelKeyMaxLength {
				return validation.Errors{
					"labels": util.ErrValidationLabelKeyLength,
				}
			}

			// Value validation
			if len(value) > util.LabelValueMaxLength {
				return validation.Errors{
					"labels": util.ErrValidationLabelValueLength,
				}
			}
		}
	}

	return nil
}

// APIExpectedMachineUpdateRequest is the data structure to capture user request to update an ExpectedMachine
type APIExpectedMachineUpdateRequest struct {
	// ID is required for batch updates (must be empty or match path value for single update)
	ID *string `json:"id"`
	// BmcMacAddress is the MAC address of the expected machine's BMC
	BmcMacAddress *string `json:"bmcMacAddress"`
	// BmcUsername is the username of the expected machine's BMC
	DefaultBmcUsername *string `json:"defaultBmcUsername"`
	// DefaultBmcPassword is the password of the expected machine's BMC
	DefaultBmcPassword *string `json:"defaultBmcPassword"`
	// ChassisSerialNumber is the serial number of the expected machine's chassis
	ChassisSerialNumber *string `json:"chassisSerialNumber"`
	// FallbackDPUSerialNumbers is the serial numbers of the expected machine's fallback DPUs
	FallbackDPUSerialNumbers []string `json:"fallbackDPUSerialNumbers"`
	// SkuId is the optional UUID for an SKU
	SkuID *string `json:"skuId"`
	// Labels is the labels of the expected machine
	Labels map[string]string `json:"labels"`
}

// Validate ensure the values passed in request are acceptable
func (emur *APIExpectedMachineUpdateRequest) Validate() error {
	if emur.ID != nil {
		if *emur.ID == "" {
			return validation.Errors{
				"id": errors.New("ID cannot be empty"),
			}
		}
		if _, err := uuid.Parse(*emur.ID); err != nil {
			return validation.Errors{
				"id": errors.New("ID must be a valid UUID"),
			}
		}
	}

	err := validation.ValidateStruct(emur,
		validation.Field(&emur.DefaultBmcUsername,
			validation.NilOrNotEmpty.Error("BMC Username cannot be empty"),
			validation.Length(1, 16).Error("BMC Username must be 1-16 characters")),
		validation.Field(&emur.DefaultBmcPassword,
			validation.NilOrNotEmpty.Error("BMC Password cannot be empty"),
			validation.Length(1, 20).Error("BMC Password must be 1-20 characters")),
		validation.Field(&emur.ChassisSerialNumber,
			validation.NilOrNotEmpty.Error("Chassis Serial Number cannot be empty"),
			validation.Length(1, 32).Error("Chassis Serial Number must be 1-32 characters")),
		validation.Field(&emur.SkuID,
			validation.NilOrNotEmpty.Error("SkuID cannot be empty")),
	)

	if err != nil {
		return err
	}

	// Labels validation
	if emur.Labels != nil {
		if len(emur.Labels) > util.LabelCountMax {
			return validation.Errors{
				"labels": util.ErrValidationLabelCount,
			}
		}

		for key, value := range emur.Labels {
			if key == "" {
				return validation.Errors{
					"labels": util.ErrValidationLabelKeyEmpty,
				}
			}

			// Key validation
			if len(key) > util.LabelKeyMaxLength {
				return validation.Errors{
					"labels": util.ErrValidationLabelKeyLength,
				}
			}

			// Value validation
			if len(value) > util.LabelValueMaxLength {
				return validation.Errors{
					"labels": util.ErrValidationLabelValueLength,
				}
			}
		}
	}

	return nil
}

// APIExpectedMachine is the data structure to capture API representation of an ExpectedMachine
type APIExpectedMachine struct {
	// ID is the ID of this Expected Machine
	ID uuid.UUID `json:"id"`
	// BmcMacAddress is the MAC address of the expected machine's BMC
	BmcMacAddress string `json:"bmcMacAddress"`
	// SiteID is the ID of the site this machine belongs to
	SiteID uuid.UUID `json:"siteId"`
	// Site is the site information
	Site *APISite `json:"site,omitempty"`
	// ChassisSerialNumber is the serial number of the expected machine's chassis
	ChassisSerialNumber string `json:"chassisSerialNumber"`
	// FallbackDPUSerialNumbers is the serial numbers of the expected machine's fallback DPUs
	FallbackDPUSerialNumbers []string `json:"fallbackDPUSerialNumbers"`
	// SkuID is the ID of the SKU
	SkuID *string `json:"skuId"`
	// Sku is the SKU information
	Sku *APISku `json:"sku,omitempty"`
	// MachineID is the ID of the Machine associated with this Expected Machine
	MachineID *string `json:"machineId"`
	// Machine is the optional Machine information associated with this Expected Machine
	Machine *APIMachineSummary `json:"machine,omitempty"`
	// Labels is the labels of the expected machine
	Labels map[string]string `json:"labels"`
	// Created indicates the ISO datetime string for when the ExpectedMachine was created
	Created time.Time `json:"created"`
	// Updated indicates the ISO datetime string for when the ExpectedMachine was last updated
	Updated time.Time `json:"updated"`
}

// NewAPIExpectedMachine accepts a DB layer ExpectedMachine object and returns an API object
func NewAPIExpectedMachine(dibp *cdbm.ExpectedMachine) *APIExpectedMachine {
	apiem := &APIExpectedMachine{
		ID:                       dibp.ID,
		BmcMacAddress:            dibp.BmcMacAddress,
		SiteID:                   dibp.SiteID,
		ChassisSerialNumber:      dibp.ChassisSerialNumber,
		FallbackDPUSerialNumbers: dibp.FallbackDpuSerialNumbers,
		SkuID:                    dibp.SkuID,
		MachineID:                dibp.MachineID,
		Labels:                   dibp.Labels,
		Created:                  dibp.Created,
		Updated:                  dibp.Updated,
	}

	// Expand Site details if available
	if dibp.Site != nil {
		site := NewAPISite(*dibp.Site, []cdbm.StatusDetail{}, nil)
		apiem.Site = &site
	}

	// Expand SKU details if available
	if dibp.Sku != nil {
		sku := NewAPISku(dibp.Sku)
		apiem.Sku = sku
	}

	// Expand Machine details if available
	if dibp.Machine != nil {
		machine := NewAPIMachineSummary(dibp.Machine)
		apiem.Machine = machine
	}

	return apiem
}
