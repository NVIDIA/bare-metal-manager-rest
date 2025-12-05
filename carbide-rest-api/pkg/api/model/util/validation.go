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


package util

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/google/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	LeadingWhitespaceRegexp  = regexp.MustCompile("^\\s+.*")
	TrailingWhitespaceRegexp = regexp.MustCompile(".*\\s+$")
	NotAllWhitespaceRegexp   = regexp.MustCompile("[^\\s]+")
	ShaHashRegex             = regexp.MustCompile("^[A-Fa-f0-9]+$")
	DiskImagePathRegex       = regexp.MustCompile("^/dev/(:?nvme\\d+n\\d+|sd*)")

	ValidationErrorExpectStringType          = errors.New("expect name field to be string type")
	ValidationErrorNameHasLeadingWhitespace  = errors.New("name field has leading whitespace")
	ValidationErrorNameHasTrailingWhitespace = errors.New("name field has trailing whitespace")
	ValidationErrorNameFieldIsEmpty          = errors.New("name field is empty")
	ValidationErrorOnlyWhitespace            = errors.New("field consists only of whitespace")

	// Label restrictions
	LabelKeyMaxLength   = 255
	LabelValueMaxLength = 255
	LabelCountMax       = 10

	// Label validation error messages
	ErrValidationLabelKeyEmpty    = errors.New("one or more labels do not have a key specified")
	ErrValidationLabelKeyLength   = fmt.Errorf("label key must contain at least 1 character and a maximum of %v characters", LabelKeyMaxLength)
	ErrValidationLabelValueLength = fmt.Errorf("label value cannot exceed a maximum of %v characters", LabelValueMaxLength)
	ErrValidationLabelCount       = fmt.Errorf("up to %v key/value pairs can be specified in labels", LabelCountMax)
)

// util.GetUUIDPtrToStrPtr is a utility function to return string pointer from uuid pointer
func GetUUIDPtrToStrPtr(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	s := id.String()
	return &s
}

// ValidateNested is a utility function to validate nested struct
func ValidateNested(target interface{}, fieldRules ...*validation.FieldRules) *validation.FieldRules {
	if target == nil {
		return nil
	}

	return validation.Field(target, validation.By(func(value interface{}) error {
		valueV := reflect.Indirect(reflect.ValueOf(value))
		if valueV.CanAddr() {
			addr := valueV.Addr().Interface()
			return validation.ValidateStruct(addr, fieldRules...)
		}
		return validation.ValidateStruct(target, fieldRules...)
	}))
}

// ValidateNameCharacters is a utility function to lexically validate the name field
// currently checks for leading or trailing whitespaces
// could evolve in the future, if there are more constraints
func ValidateNameCharacters(value interface{}) error {
	s, ok := value.(string)
	var name string
	if !ok {
		// check for string pointer
		sPtr, ok := value.(*string)
		if !ok || sPtr == nil {
			return ValidationErrorExpectStringType
		}
		name = *sPtr
	} else {
		name = s
	}
	if LeadingWhitespaceRegexp.Match([]byte(name)) {
		return ValidationErrorNameHasLeadingWhitespace
	}
	if TrailingWhitespaceRegexp.Match([]byte(name)) {
		return ValidationErrorNameHasTrailingWhitespace
	}
	return nil
}

// IsNilOrEmptyStrPtr is a utility function to check if the string pointer is nil or the underlying value is empty
func IsNilOrEmptyStrPtr(s *string) bool {
	return s == nil || *s == ""
}

// IsEmptyStrPtr is a utility function to check if the underlying value of a string pointer is empty
func IsEmptyStrPtr(s *string) bool {
	return s != nil && *s == ""
}
