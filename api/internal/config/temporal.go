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
	"crypto/tls"
	"fmt"

	"github.com/nvidia/carbide-rest/cert-manager/pkg/core"
	cwfns "github.com/nvidia/carbide-rest/workflow/pkg/namespace"
)

// TemporalConfig holds configuration for Temporal communication
type TemporalConfig struct {
	Host          string
	Port          int
	ServerName    string
	Namespace     string
	Queue         string
	EncryptionKey string
	TLSenabled    bool
	ClientTLSCfg  *tls.Config
	dynTLS        *core.DynTLSCfg
}

// GetHostPort returns the concatenated host & port
func (tcfg *TemporalConfig) GetHostPort() string {
	return fmt.Sprintf("%v:%v", tcfg.Host, tcfg.Port)
}

func (tcfg *TemporalConfig) Close() {
	if tcfg.dynTLS != nil {
		tcfg.dynTLS.Close()
	}
}

// NewTemporalConfig initializes and returns a configuration object for managing Temporal
func NewTemporalConfig(host string, port int, serverName string, namespace string, queue string, encryptionKey string, tlsEnabled bool, certPath string, keyPath string, caPath string) (*TemporalConfig, error) {
	var dynTLS *core.DynTLSCfg
	var clientTLSCfg *tls.Config

	if tlsEnabled {
		var err error

		dynTLS, err = core.NewDynTLSCfg(keyPath, certPath, caPath)
		if err != nil {
			return nil, err
		}

		baseTLSCfg := &tls.Config{
			ServerName:         fmt.Sprintf("%s.%s", cwfns.CloudNamespace, serverName),
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false,
		}

		clientTLSCfg = dynTLS.WithTLSCfg(baseTLSCfg).ClientCfg()
	}

	return &TemporalConfig{
		Host:          host,
		Port:          port,
		Namespace:     namespace,
		Queue:         queue,
		ServerName:    serverName,
		EncryptionKey: encryptionKey,
		TLSenabled:    tlsEnabled,
		ClientTLSCfg:  clientTLSCfg,
		dynTLS:        dynTLS,
	}, nil
}
