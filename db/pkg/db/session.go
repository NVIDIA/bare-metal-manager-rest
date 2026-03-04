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
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bunotel"
)

// Session is a wrapper for an ORM DB object
type Session struct {
	DBName       string
	DB           *bun.DB
	pool         *pgxpool.Pool
	errorChecker ErrorChecker
}

// NewSession creates and returns a new session object using pgx v5 + pgxpool.
func NewSession(ctx context.Context, host string, port int, dbName string, user string, password string, caCertPath string) (*Session, error) {
	configDSN := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", url.PathEscape(user), url.PathEscape(password), host, port, dbName)

	if caCertPath != "" {
		configDSN = fmt.Sprintf("%v?sslmode=verify-full&sslrootcert=%v", configDSN, caCertPath)
	}

	pool, err := pgxpool.New(ctx, configDSN)
	if err != nil {
		return nil, err
	}

	sqldb := stdlib.OpenDBFromPool(pool)
	db := bun.NewDB(sqldb, pgdialect.New())

	// if tracing service name is configured, add otel hooks
	if os.Getenv("TRACING_SERVICE_NAME") != "" {
		db.AddQueryHook(bunotel.NewQueryHook(
			bunotel.WithDBName(dbName),
			bunotel.WithFormattedQueries(true),
		))
	}

	return &Session{
		DBName:       dbName,
		DB:           db,
		pool:         pool,
		errorChecker: &PostgresErrorChecker{},
	}, nil
}

// NewSessionFromConfig creates a Session from a Config.
func NewSessionFromConfig(ctx context.Context, c Config) (*Session, error) {
	pool, err := pgxpool.New(ctx, c.BuildDSN())
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New())

	// if tracing service name is configured, add otel hooks
	if os.Getenv("TRACING_SERVICE_NAME") != "" {
		db.AddQueryHook(bunotel.NewQueryHook(
			bunotel.WithDBName(c.DBName),
			bunotel.WithFormattedQueries(true),
		))
	}

	return &Session{
		DBName:       c.DBName,
		DB:           db,
		pool:         pool,
		errorChecker: &PostgresErrorChecker{},
	}, nil
}

// Close closes the session and the underlying connection pool.
func (s *Session) Close() {
	s.DB.Close()

	if s.pool != nil {
		s.pool.Close()
	}
}

// GetErrorChecker returns the error classifier for this session.
func (s *Session) GetErrorChecker() ErrorChecker { //nolint:ireturn // interface by design
	return s.errorChecker
}

// BeginTx begins a new transaction with default options.
func (s *Session) BeginTx(ctx context.Context) (bun.Tx, error) {
	return s.DB.BeginTx(ctx, &sql.TxOptions{}) //nolint:exhaustruct,wrapcheck // default options; thin wrapper
}

// RunInTx executes a function within a transaction.
func (s *Session) RunInTx(
	ctx context.Context,
	fn func(ctx context.Context, tx bun.Tx) error,
) error {
	return s.DB.RunInTx(ctx, &sql.TxOptions{}, fn) //nolint:exhaustruct,wrapcheck // default options; thin wrapper
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
