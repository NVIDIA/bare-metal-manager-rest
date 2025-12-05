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
	"github.com/stretchr/testify/assert"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	"testing"
)

func TestAPIMachineValidationTestCreateRequest_Validate(t *testing.T) {
	tests := []struct {
		desc      string
		obj       APIMachineValidationTestCreateRequest
		expectErr bool
	}{
		{
			desc:      "no error",
			obj:       APIMachineValidationTestCreateRequest{Name: "test-1", Command: "/bin/sh/test1", Args: "-p 12"},
			expectErr: false,
		},
		{
			desc:      "error no Name",
			obj:       APIMachineValidationTestCreateRequest{Command: "/bin/sh/test1", Args: "-p 12"},
			expectErr: true,
		},
		{
			desc:      "error no Command",
			obj:       APIMachineValidationTestCreateRequest{Name: "test-1", Args: "-p 12"},
			expectErr: true,
		},
		{
			desc:      "error no args",
			obj:       APIMachineValidationTestCreateRequest{Name: "test-1", Command: "/bin/sh/test1"},
			expectErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.obj.Validate()
			assert.Equal(t, tc.expectErr, err != nil)
			if err != nil {
				fmt.Println(err.Error())
			}
		})
	}
}

func TestAPIMachineValidationExternalConfigCreateRequest_Validate(t *testing.T) {
	tests := []struct {
		desc      string
		obj       APIMachineValidationExternalConfigCreateRequest
		expectErr bool
	}{
		{
			desc:      "no error",
			obj:       APIMachineValidationExternalConfigCreateRequest{Name: "test-1", Description: cdb.GetStrPtr("test description"), Config: []byte{0, 1, 12}},
			expectErr: false,
		},
		{
			desc:      "no error with no description",
			obj:       APIMachineValidationExternalConfigCreateRequest{Name: "test-1", Config: []byte{0, 1, 12}},
			expectErr: false,
		},
		{
			desc:      "error no Name",
			obj:       APIMachineValidationExternalConfigCreateRequest{Config: []byte{0, 1, 12}},
			expectErr: true,
		},
		{
			desc:      "error no Config",
			obj:       APIMachineValidationExternalConfigCreateRequest{Name: "test-1"},
			expectErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.obj.Validate()
			assert.Equal(t, tc.expectErr, err != nil)
			if err != nil {
				fmt.Println(err.Error())
			}
		})
	}
}
