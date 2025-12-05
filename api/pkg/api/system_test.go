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


package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSystemAPIRoutes(t *testing.T) {
	tests := []struct {
		name string
		want []Route
	}{
		{
			name: "test initializing system API routes",
			want: []Route{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSystemAPIRoutes()

			assert.Equal(t, len(got), 2)
		})
	}
}
