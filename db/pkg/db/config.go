/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package db

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/nvidia/bare-metal-manager-rest/common/pkg/credential"
)

// Config represents the configuration needed to connect to a database.
type Config struct {
	Host              string
	Port              int
	DBName            string
	Credential        credential.Credential
	CACertificatePath string
}

// Validate checks if the Config fields are set correctly.
func (c *Config) Validate() error {
	if c.Host == "" {
		return errors.New("host is required")
	}

	if c.Port <= 0 || c.Port > 65535 {
		return errors.New("port must be between (0, 65535]")
	}

	if c.DBName == "" {
		return errors.New("database name is required")
	}

	if !c.Credential.IsValid() {
		return errors.New("valid credential is required")
	}

	return nil
}

// ConfigFromEnv builds a Config from environment variables.
// Reads: DB_ADDR or DB_HOST (host), DB_PORT (port), DB_USER, DB_PASSWORD,
// DB_DATABASE or DB_NAME (database name), DB_CERT_PATH (optional CA certificate).
func ConfigFromEnv() (Config, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return Config{}, ErrInvalidPort
	}

	cred := credential.NewFromEnv("DB_USER", "DB_PASSWORD")
	if !cred.IsValid() {
		return Config{}, ErrInvalidCredential
	}

	host := os.Getenv("DB_ADDR")
	if host == "" {
		host = os.Getenv("DB_HOST")
	}

	dbName := os.Getenv("DB_DATABASE")
	if dbName == "" {
		dbName = os.Getenv("DB_NAME")
	}

	return Config{
		Host:              host,
		Port:              port,
		Credential:        cred,
		DBName:            dbName,
		CACertificatePath: os.Getenv("DB_CERT_PATH"),
	}, nil
}

// BuildDSN builds the Data Source Name (DSN) string for connecting to
// the database.
func (c *Config) BuildDSN() string {
	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=",
		url.PathEscape(c.Credential.User),
		url.PathEscape(c.Credential.Password.Value),
		c.Host,
		c.Port,
		c.DBName,
	)

	if len(c.CACertificatePath) > 0 {
		dsn += fmt.Sprintf("prefer&sslrootcert=%v", c.CACertificatePath)
	} else {
		dsn += "disable"
	}

	return dsn
}
