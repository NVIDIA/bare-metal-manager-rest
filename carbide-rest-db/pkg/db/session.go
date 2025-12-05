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

package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bunotel"
)

// Session is a wrapper for an ORM DB object
type Session struct {
	DBName string
	DB     *bun.DB
}

// Close closes the session
func (s *Session) Close() {
	s.DB.Close()
}

// NewSession creates and returns a new session object
func NewSession(host string, port int, dbName string, user string, password string, caCertPath string) (*Session, error) {
	configDSN := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, password, host, port, dbName)

	if caCertPath != "" {
		configDSN = fmt.Sprintf("%v?sslmode=verify-full&sslrootcert=%v", configDSN, caCertPath)
	}

	config, err := pgx.ParseConfig(configDSN)
	if err != nil {
		return nil, err
	}
	config.PreferSimpleProtocol = true

	sqldb := stdlib.OpenDB(*config)
	db := bun.NewDB(sqldb, pgdialect.New())

	// Uncomment to see failing queries
	// db.AddQueryHook(bundebug.NewQueryHook())

	// Uncomment to see all queries
	// bundebug.NewQueryHook(bundebug.WithVerbose(true))

	// if tracing service name is configured, add otel hooks
	if os.Getenv("TRACING_SERVICE_NAME") != "" {
		db.AddQueryHook(bunotel.NewQueryHook(
			bunotel.WithDBName(dbName),
			bunotel.WithFormattedQueries(true),
		))
	}

	return &Session{
		DBName: dbName,
		DB:     db,
	}, nil
}

// acquireAdvisoryLock will "try" to take the specified advisory lock
// on the session - which is a specific connection to the DB. the
// advisory lock has to be released from that same connection
//
// However, There are 2 problems with session advisory locks. and
// hence session advisory locks are not recommended. and transaction
// advisory locks are recommended. hence this method is not exposed.
//
// Problem 1.
// Since the application only works with a connection
// pool implemented in database/sql, the specific connection on which the
// session advisory lock is acquired is not visible to the application.
// Hence the application is unable to correctly unlock the lock,
// since it does not have the connection handle (which was used to lock)
//
// Problem 2.
// session advisory locks can be taken multiple times on the same
// connection. It is possible therefore that 2 api's could both acquire the
// session lock (since they both use the same connection pool, and it could
// so happen that they used the same connection to execute their sqls).
func (s *Session) acquireAdvisoryLock(ctx context.Context, lockID uint64) error {
	query := fmt.Sprintf("pg_try_advisory_lock(%d)", lockID)
	value := false
	err := s.DB.NewSelect().ColumnExpr(query).Scan(ctx, &value)
	if err != nil {
		return err
	}
	if !value {
		return ErrSessionAdvisoryLockFailed
	}
	return nil
}

// releaseAdvisoryLock will release the advisory lock taken earlier
func (s *Session) releaseAdvisoryLock(ctx context.Context, lockID uint64) error {
	query := fmt.Sprintf("pg_advisory_unlock(%d)", lockID)
	value := false
	err := s.DB.NewSelect().ColumnExpr(query).Scan(ctx, &value)
	if err != nil {
		return err
	}
	if !value {
		return ErrSessionAdvisoryLockUnlockFailed
	}
	return nil
}
