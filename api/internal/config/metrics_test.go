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


package config

import (
	"testing"
)

func TestMetricsConfig(t *testing.T) {
	type args struct {
		enabled  bool
		port     int
	}

	mcfg := MetricsConfig{
		Enabled:  true,
		Port:     6930,
	}

	tests := []struct {
		name string
		args args
		want *MetricsConfig
	}{
		{
			name: "initialize Metrics config",
			args: args{
				enabled:  true,
				port:     mcfg.Port,
			},
			want: &mcfg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMetricsConfig(tt.args.enabled, tt.args.port)

			if p := got.Port; p != tt.want.Port {
				t.Errorf("got.Port = %v, want %v", p, tt.want.Port)
			}

			if got := got.GetListenAddr(); got != tt.want.GetListenAddr() {
				t.Errorf("GetListenAddr() = %v, want %v", got, tt.want.GetListenAddr())
			}
		})
	}
}
