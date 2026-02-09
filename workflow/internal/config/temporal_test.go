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
	"crypto/tls"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/nvidia/carbide-rest/cert-manager/pkg/core"
)

func TestNewTemporalConfig(t *testing.T) {
	type args struct {
		host          string
		port          int
		serverName    string
		namespace     string
		queue         string
		encryptionKey string
	}

	keyPath, certPath := SetupTestCerts(t)
	defer os.Remove(keyPath)
	defer os.Remove(certPath)
	d, err := core.NewDynTLSCfg(keyPath, certPath, certPath)
	assert.NoError(t, err)
	defer d.Close()

	tcfg := TemporalConfig{
		Host:          "localhost",
		Port:          7233,
		ServerName:    "temporal.local",
		Namespace:     "cloud",
		Queue:         "cloud",
		EncryptionKey: "test",
		ClientTLSCfg: &tls.Config{
			ServerName:         fmt.Sprintf("%s.%s", "cloud", "temporal.local"),
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false,
		},
		dynTLS: d,
	}

	tests := []struct {
		name string
		args args
		want *TemporalConfig
	}{
		{
			name: "initialize Temporal config",
			args: args{
				host:          tcfg.Host,
				port:          tcfg.Port,
				serverName:    tcfg.ServerName,
				namespace:     tcfg.Namespace,
				queue:         tcfg.Queue,
				encryptionKey: tcfg.EncryptionKey,
			},
			want: &tcfg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTemporalConfig(tt.args.host, tt.args.port, tt.args.serverName, tt.args.namespace, tt.args.queue, tt.args.encryptionKey, true, certPath, keyPath, certPath)
			assert.NoError(t, err)
			defer got.Close()

			if sn := got.ServerName; sn != tt.want.ServerName {
				t.Errorf("got.ServerName = %v, want %v", sn, tt.want.ServerName)
			}
			if got := got.GetHostPort(); got != tt.want.GetHostPort() {
				t.Errorf("GetHostPort() = %v, want %v", got, tt.want.GetHostPort())
			}
		})
	}
}
