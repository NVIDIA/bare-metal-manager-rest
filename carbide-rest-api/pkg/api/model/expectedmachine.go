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

	"github.com/google/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	validationIs "github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/NVIDIA/carbide-rest-api/carbide-rest-api/pkg/api/model/util"
	cdbm "github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/db/model"
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
	SkuId *string `json:"skuId"`
	// Labels is the labels of the expected machine
	Labels map[string]string `json:"labels"`
}

// Validate ensure the values passed in request are acceptable
func (emcr *APIExpectedMachineCreateRequest) Validate() error {
	err := validation.ValidateStruct(emcr,
		validation.Field(&emcr.BmcMacAddress,
			validation.Required.Error(validationErrorValueRequired),
			validationIs.MAC),
		validation.Field(&emcr.DefaultBmcUsername,
			validation.Length(0, 16).Error("BMC username must be 16 characters or less")),
		validation.Field(&emcr.DefaultBmcPassword,
			validation.Length(0, 20).Error("BMC password must be 20 characters or less")),
		validation.Field(&emcr.ChassisSerialNumber,
			validation.Required.Error(validationErrorValueRequired),
			validation.Length(1, 32).Error("Chassis serial number must be 32 characters or less")),
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
// For now same as CreateRequest
type APIExpectedMachineUpdateRequest struct {
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
	SkuId *string `json:"skuId"`
	// Labels is the labels of the expected machine
	Labels map[string]string `json:"labels"`
}

// Validate ensure the values passed in request are acceptable
func (emur *APIExpectedMachineUpdateRequest) Validate() error {
	// Validate DefaultBmcUsername: if provided, cannot be empty
	if emur.DefaultBmcUsername != nil && *emur.DefaultBmcUsername == "" {
		return validation.Errors{
			"defaultBmcUsername": validation.NewError("validation_length_out_of_range", "BMC Username cannot be empty"),
		}
	}

	// Validate DefaultBmcPassword: if provided, cannot be empty
	if emur.DefaultBmcPassword != nil && *emur.DefaultBmcPassword == "" {
		return validation.Errors{
			"defaultBmcPassword": validation.NewError("validation_length_out_of_range", "BMC Password cannot be empty"),
		}
	}

	// Validate ChassisSerialNumber: if provided, must be 1-32 characters
	if emur.ChassisSerialNumber != nil && *emur.ChassisSerialNumber == "" {
		return validation.Errors{
			"chassisSerialNumber": validation.NewError("validation_length_out_of_range", "Chassis Serial Number number must be 1-32 characters"),
		}
	}

	err := validation.ValidateStruct(emur,
		validation.Field(&emur.DefaultBmcUsername,
			validation.Length(1, 16).Error("BMC Username must be 1-16 characters")),
		validation.Field(&emur.DefaultBmcPassword,
			validation.Length(1, 20).Error("BMC Password must be 1-20 characters")),
		validation.Field(&emur.ChassisSerialNumber,
			validation.Length(1, 32).Error("Chassis Serial Number must be 1-32 characters")),
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
	Site *APISite `json:"site"`
	// ChassisSerialNumber is the serial number of the expected machine's chassis
	ChassisSerialNumber string `json:"chassisSerialNumber"`
	// FallbackDPUSerialNumbers is the serial numbers of the expected machine's fallback DPUs
	FallbackDPUSerialNumbers []string `json:"fallbackDPUSerialNumbers"`
	// SkuID is the ID of the SKU
	SkuID *string `json:"skuId"`
	// Sku is the SKU information
	Sku *APISku `json:"sku"`
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
		Sku:                      NewAPISku(dibp.Sku),
		Labels:                   dibp.Labels,
		Created:                  dibp.Created,
		Updated:                  dibp.Updated,
	}
	// Map Site if available
	if dibp.Site != nil {
		site := NewAPISite(*dibp.Site, []cdbm.StatusDetail{}, nil)
		apiem.Site = &site
	}

	return apiem
}
