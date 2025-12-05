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
	"fmt"
)

// TemporalConfig holds configuration for Temporal communication
type MetricsConfig struct {
	Enabled       bool
	Port          int
}

// GetListenAddr returns the local address for listen socket.
func (mcfg *MetricsConfig) GetListenAddr() string {
	return fmt.Sprintf(":%v", mcfg.Port)
}

// NewMetricsConfig initializes and returns a configuration object for managing Metrics
func NewMetricsConfig(enabled bool, port int) *MetricsConfig {
	return &MetricsConfig{
		Enabled:       enabled,
		Port:          port,
	}
}

