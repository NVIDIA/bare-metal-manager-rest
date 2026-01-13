// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

// Package config provides configuration for site-workflow activities and workflows.
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

// Default timeout values
const (
	// DefaultInventoryActivityTimeout is the default timeout for inventory collection activities.
	// This needs to be long enough to handle large machine inventories (hundreds or thousands of machines).
	DefaultInventoryActivityTimeout = 5 * time.Minute

	// DefaultActivityTimeout is the default timeout for standard activities.
	DefaultActivityTimeout = 2 * time.Minute
)

// Config holds the site-workflow configuration values.
// These are initialized at startup from environment variables.
type Config struct {
	// InventoryActivityTimeout is the StartToCloseTimeout for inventory collection activities.
	// Set via INVENTORY_ACTIVITY_TIMEOUT_MINUTES environment variable.
	// Default: 5 minutes
	InventoryActivityTimeout time.Duration
}

var globalConfig *Config

// init loads configuration from environment variables at startup.
func init() {
	globalConfig = loadConfig()
}

// loadConfig reads configuration from environment variables with defaults.
func loadConfig() *Config {
	cfg := &Config{
		InventoryActivityTimeout: DefaultInventoryActivityTimeout,
	}

	// INVENTORY_ACTIVITY_TIMEOUT_MINUTES: timeout in minutes for inventory collection activities
	if envVal := os.Getenv("INVENTORY_ACTIVITY_TIMEOUT_MINUTES"); envVal != "" {
		if minutes, err := strconv.Atoi(envVal); err == nil && minutes > 0 {
			cfg.InventoryActivityTimeout = time.Duration(minutes) * time.Minute
			log.Info().
				Int("timeout_minutes", minutes).
				Msg("Config: Loaded inventory activity timeout from environment")
		} else {
			log.Warn().
				Str("value", envVal).
				Msg("Config: Invalid INVENTORY_ACTIVITY_TIMEOUT_MINUTES, using default")
		}
	}

	return cfg
}

// GetInventoryActivityTimeout returns the configured timeout for inventory collection activities.
func GetInventoryActivityTimeout() time.Duration {
	if globalConfig == nil {
		return DefaultInventoryActivityTimeout
	}
	return globalConfig.InventoryActivityTimeout
}

// GetDefaultActivityTimeout returns the default timeout for standard activities.
func GetDefaultActivityTimeout() time.Duration {
	return DefaultActivityTimeout
}
