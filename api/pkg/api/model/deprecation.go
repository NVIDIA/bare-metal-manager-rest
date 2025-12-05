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
	"fmt"
	"time"
)

const (
	deprecationPreTemplate  = `"'%s' is being deprecated%s. Please take action prior to the specified date"`
	deprecationPostTemplate = `"'%s' has been deprecated%s. Please take action immediately"`

	// DeprecationTypeAttribute denotes a deprecation of an API model attribute
	DeprecationTypeAttribute = "Attribute"
	// DeprecationTypeQueryParam denotes a deprecation of an API query parameter
	DeprecationTypeQueryParam = "QueryParam"
	// DeprecationTypeEndpoint denotes a deprecation of an API endpoint
	DeprecationTypeEndpoint = "Endpoint"
)

// DeprecatdEntity denotes an entity that is being deprecated
type DeprecatedEntity struct {
	OldValue     string
	NewValue     *string
	Type         string
	TakeActionBy time.Time
}

// APIDeprecation captures API representation of a deprecation message
type APIDeprecation struct {
	// Field denotes the field that is deprecated (optional)
	Attribute *string `json:"attribute,omitempty"`
	// Field denotes the field that is deprecated (optional)
	QueryParam *string `json:"queryparam,omitempty"`
	// Endpoint denotes the endpoint that is deprecated (optional)
	Endpoint *string `json:"endpoint,omitempty"`
	// ReplacedBy denotes the field that replaces the deprecated field (optional)
	ReplacedBy *string `json:"replacedby"`
	// Effective indicates the ISO datetime string for when the deprecation takes effect
	TakeActionBy time.Time `json:"effective"`
	// Notice describes the deprecated field
	Notice string `json:"notice"`
}

// NewAPIDeprecation creates an API deprecation object from parameters
func NewAPIDeprecation(de DeprecatedEntity) APIDeprecation {
	apiDeprecation := APIDeprecation{
		TakeActionBy: de.TakeActionBy,
	}

	if de.Type == DeprecationTypeAttribute {
		apiDeprecation.Attribute = &de.OldValue
	} else if de.Type == DeprecationTypeQueryParam {
		apiDeprecation.QueryParam = &de.OldValue
	} else if de.Type == DeprecationTypeEndpoint {
		apiDeprecation.Endpoint = &de.OldValue
	}

	if de.NewValue != nil {
		apiDeprecation.ReplacedBy = de.NewValue
	}

	noticeReplacedBy := ""
	if de.NewValue != nil {
		noticeReplacedBy = fmt.Sprintf(" in favor of '%s'", *de.NewValue)
	}

	if de.TakeActionBy.After(time.Now()) {
		apiDeprecation.Notice = fmt.Sprintf(deprecationPreTemplate, de.OldValue, noticeReplacedBy)
	} else {
		apiDeprecation.Notice = fmt.Sprintf(deprecationPostTemplate, de.OldValue, noticeReplacedBy)
	}

	return apiDeprecation
}
