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


package site

import (
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/nvidia/carbide-rest/api/internal/config"

	temporalClient "go.temporal.io/sdk/client"
)

func TestNewClientPool(t *testing.T) {
	type args struct {
		tcfg *config.TemporalConfig
	}

	keyPath, certPath := config.SetupTestCerts(t)
	defer os.Remove(keyPath)
	defer os.Remove(certPath)

	cfg := config.NewConfig()
	cfg.SetTemporalCertPath(certPath)
	cfg.SetTemporalKeyPath(keyPath)
	cfg.SetTemporalCaPath(certPath)

	tcfg, err := cfg.GetTemporalConfig()
	assert.NoError(t, err)

	tests := []struct {
		name string
		args args
		want *ClientPool
	}{
		{
			name: "test Site client pool initializer",
			args: args{
				tcfg: tcfg,
			},
			want: &ClientPool{
				tcfg:        tcfg,
				IDClientMap: map[string]temporalClient.Client{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClientPool(tt.args.tcfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSitePool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientPool_GetClientByID(t *testing.T) {
	type fields struct {
		tcfg *config.TemporalConfig
	}
	type args struct {
		siteID uuid.UUID
	}

	keyPath, certPath := config.SetupTestCerts(t)
	defer os.Remove(keyPath)
	defer os.Remove(certPath)

	cfg := config.NewConfig()
	cfg.SetTemporalCertPath(certPath)
	cfg.SetTemporalKeyPath(keyPath)
	cfg.SetTemporalCaPath(certPath)

	tcfg, err := cfg.GetTemporalConfig()
	assert.NoError(t, err)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    temporalClient.Client
		wantErr bool
	}{
		{
			name: "test retrieving client for given site ID",
			fields: fields{
				tcfg: tcfg,
			},
			args: args{
				siteID: uuid.New(),
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := NewClientPool(tt.fields.tcfg)
			_, err := cp.GetClientByID(tt.args.siteID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientPool.GetClientByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
