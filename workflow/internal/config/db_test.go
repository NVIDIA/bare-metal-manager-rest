// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package config

import (
	"reflect"
	"testing"
)

func TestNewDBConfig(t *testing.T) {
	type args struct {
		host     string
		port     int
		name     string
		user     string
		password string
	}

	dbcfg := DBConfig{
		Host:     "localhost",
		Port:     5432,
		Name:     "forge",
		User:     "forge",
		Password: "test123",
	}

	tests := []struct {
		name string
		args args
		want *DBConfig
	}{
		{
			name: "initialize database config",
			args: args{
				host:     dbcfg.Host,
				port:     dbcfg.Port,
				name:     dbcfg.Name,
				user:     dbcfg.User,
				password: dbcfg.Password,
			},
			want: &dbcfg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDBConfig(tt.args.host, tt.args.port, tt.args.name, tt.args.user, tt.args.password)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDBConfig() = %v, want %v", got, tt.want)
			}

			if got := got.GetHostPort(); got != tt.want.GetHostPort() {
				t.Errorf("GetHostPort() = %v, want %v", got, tt.want.GetHostPort())
			}
		})
	}
}
