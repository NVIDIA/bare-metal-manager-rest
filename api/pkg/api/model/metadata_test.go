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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/nvidia/carbide-rest/api/pkg/metadata"
)

func TestNewAPIMetadata(t *testing.T) {
	tests := []struct {
		name string
		want *APIMetadata
	}{
		{
			name: "test initializing API model for HealthCheck",
			want: &APIMetadata{
				Version:   metadata.Version,
				BuildTime: metadata.BuildTime,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAPIMetadata()

			assert.Equal(t, tt.want.Version, got.Version)
			assert.Equal(t, tt.want.BuildTime, got.BuildTime)
		})
	}
}
