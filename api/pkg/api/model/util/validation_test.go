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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateNameCharacters(t *testing.T) {
	val := 0
	// test error when string not passed
	assert.NotNil(t, ValidateNameCharacters(val))
	assert.NotNil(t, ValidateNameCharacters(&val))
	assert.NotNil(t, ValidateNameCharacters(nil))
	tests := []struct {
		desc      string
		names     []string
		expectErr bool
	}{
		{
			desc:      "error with leading whitespaces",
			names:     []string{" hello", "\thello", "\nhello", "     "},
			expectErr: true,
		},
		{
			desc:      "errors with trailing whitespaces",
			names:     []string{"hello ", "hello\t", "hello\n"},
			expectErr: true,
		},
		{
			desc:      "success cases",
			names:     []string{"hel lo", "hel \t lo", "hel&&lo"},
			expectErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			for _, s := range tc.names {
				err := ValidateNameCharacters(s)
				assert.Equal(t, tc.expectErr, err != nil)
				err = ValidateNameCharacters(&s)
				assert.Equal(t, tc.expectErr, err != nil)
			}
		})
	}
}
