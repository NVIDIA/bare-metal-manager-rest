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
	"reflect"
	"testing"
)

func TestNewAPIHealthCheck(t *testing.T) {
	type args struct {
		isHealthy    bool
		errorMessage *string
	}
	tests := []struct {
		name string
		args args
		want *APIHealthCheck
	}{
		{
			name: "test initializing API model for HealthCheck",
			args: args{
				isHealthy:    true,
				errorMessage: nil,
			},
			want: &APIHealthCheck{
				IsHealthy: true,
				Error:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPIHealthCheck(tt.args.isHealthy, tt.args.errorMessage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIHealthCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
