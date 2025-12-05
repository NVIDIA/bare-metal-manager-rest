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

import "fmt"

// DBConfig holds configuration for database access
type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

// GetHostPort returns the concatenated host & port
func (dbcfg *DBConfig) GetHostPort() string {
	return fmt.Sprintf("%v:%v", dbcfg.Host, dbcfg.Port)
}

// NewDBConfig initializes and returns a configuration object for managing database access
func NewDBConfig(host string, port int, name string, user string, password string) *DBConfig {
	return &DBConfig{
		Host:     host,
		Port:     port,
		Name:     name,
		User:     user,
		Password: password,
	}
}
