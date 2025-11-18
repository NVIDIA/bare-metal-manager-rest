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

package util

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/db"
	"github.com/uptrace/bun/extra/bundebug"

	cipam "github.com/NVIDIA/carbide-rest-api/carbide-rest-ipam"
)

// TestDBConfig describes a test DB config params
type TestDBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

// getTestDBParams returns the DB params for a test DB
func getTestDBParams() TestDBConfig {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	
	port := 30432
	if portEnv := os.Getenv("DB_PORT"); portEnv != "" {
		fmt.Sscanf(portEnv, "%d", &port)
	}
	
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "forgetest"
	}
	
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	tdbcfg := TestDBConfig{
		Host:     host,
		Port:     port,
		Name:     name,
		User:     user,
		Password: password,
	}

	return tdbcfg
}

// GetTestIpamDB returns a test IPAM DB
func GetTestIpamDB(t *testing.T) cipam.Storage {
	tdbcfg := getTestDBParams()

	ipamDB, err := cipam.NewPostgresStorage(tdbcfg.Host, fmt.Sprintf("%d", tdbcfg.Port), tdbcfg.User, tdbcfg.Password, tdbcfg.Name, cipam.SSLModeDisable)
	if err != nil {
		t.Fatal(err)
	}

	return ipamDB
}

// GetTestDBSession returns a test DB session
func GetTestDBSession(t *testing.T, reset bool) *db.Session {
	// Create test DB
	tdbcfg := getTestDBParams()

	dbSession, err := db.NewSession(tdbcfg.Host, tdbcfg.Port, "postgres", tdbcfg.User, tdbcfg.Password, "")
	if err != nil {
		t.Fatal(err)
	}

	count, err := dbSession.DB.NewSelect().Table("pg_database").Where("datname = ?", tdbcfg.Name).Count(context.Background())
	if err != nil {
		dbSession.Close()
		t.Fatal(err)
	}

	if count > 0 && reset {
		_, err = dbSession.DB.Exec("DROP DATABASE " + tdbcfg.Name)
		if err != nil {
			dbSession.Close()
			t.Fatal(err)
		}
	}

	if count == 0 || reset {
		_, err = dbSession.DB.Exec("CREATE DATABASE " + tdbcfg.Name)
		if err != nil {
			dbSession.Close()
			t.Fatal(err)
		}
	}
	// close this session
	dbSession.Close()

	// Create another session to the forgetest database
	dbSession, err = db.NewSession(tdbcfg.Host, tdbcfg.Port, tdbcfg.Name, tdbcfg.User, tdbcfg.Password, "")
	if err != nil {
		t.Fatal(err)
	}

	_, err = dbSession.DB.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm")
	if err != nil {
		dbSession.Close()
		t.Fatal(err)
	}

	if testing.Verbose() {
		connCount, err := GetDBConnectionCount(dbSession)
		if err == nil {
			fmt.Printf("connections count = %d\n", connCount)
		}
	}

	return dbSession
}

// GetDBConnectionCount returns the count of rows in the pg_stat_activity table
func GetDBConnectionCount(dbSession *db.Session) (int, error) {
	return dbSession.DB.NewSelect().Table("pg_stat_activity").Count(context.Background())
}

// TestInitDB initializes a test DB session
func TestInitDB(t *testing.T) *db.Session {
	dbSession := GetTestDBSession(t, false)
	dbSession.DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv(""),
	))
	return dbSession
}

// CleanupTestDB cleans up all test tables to ensure test isolation
// This should be called at the start of each test that modifies the database
func CleanupTestDB(t *testing.T, dbSession *db.Session) {
	t.Helper()
	ctx := context.Background()
	
	// Truncate all tables in the correct order to avoid foreign key constraint violations
	// Note: This list should match the tables used in tests
	tables := []string{
		"ssh_key_group_instance_associations",
		"ssh_key_associations",
		"ssh_keys",
		"ssh_key_groups",
		"infiniband_interfaces",
		"infiniband_partitions",
		"interfaces",
		"machine_interfaces",
		"machine_capabilities",
		"machine_instance_types",
		"status_details",
		"instances",
		"machines",
		"subnets",
		"vpc_prefixes",
		"vpcs",
		"tenant_sites",
		"ip_blocks",
		"allocation_constraints",
		"allocations",
		"domains",
		"network_security_groups",
		"sites",
		"tenant_accounts",
		"operating_system_site_associations",
		"operating_systems",
		"instance_types",
		"users",
		"tenants",
		"infrastructure_providers",
		"audit_entries",
		"fabrics",
		"skus",
	}
	
	for _, table := range tables {
		_, err := dbSession.DB.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			// Ignore errors for tables that don't exist yet
			// This makes the helper more robust for partial schemas
			if testing.Verbose() {
				t.Logf("Warning: Failed to truncate table %s: %v", table, err)
			}
		}
	}
	
	// Also clean up the IPAM prefixes table if it exists
	_, err := dbSession.DB.ExecContext(ctx, "TRUNCATE TABLE prefixes CASCADE")
	if err != nil && testing.Verbose() {
		t.Logf("Warning: Failed to truncate prefixes table: %v", err)
	}
}
