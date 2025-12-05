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


package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
)

func tenantConfigUpMigration(ctx context.Context, db *bun.DB) error {
	// Start transactions
	tx, terr := db.BeginTx(ctx, &sql.TxOptions{})
	if terr != nil {
		handlePanic(terr, "failed to begin transaction")
	}

	// Remove NOT NULL constraint from instance_type_id in instance table
	_, err := tx.Exec("ALTER TABLE instance ALTER COLUMN instance_type_id DROP NOT NULL;")
	handleError(tx, err)

	// Ensure existing column will get an empty JSON as default value
	_, err = tx.Exec("ALTER TABLE tenant ALTER COLUMN config SET DEFAULT '{}'::jsonb")
	handleError(tx, err)

	// Set the config column in tenant table to {}::jsonb
	_, err = tx.Exec("UPDATE tenant SET config='{}'::jsonb WHERE config IS NULL")
	handleError(tx, err)

	// Set the config column in tenant table to not null
	_, err = tx.Exec("ALTER TABLE tenant ALTER COLUMN config SET NOT NULL")
	handleError(tx, err)

	terr = tx.Commit()
	if terr != nil {
		handlePanic(terr, "failed to commit transaction")
	}

	fmt.Print(" [up migration] ")
	return nil
}

func init() {
	Migrations.MustRegister(tenantConfigUpMigration, func(_ context.Context, _ *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
