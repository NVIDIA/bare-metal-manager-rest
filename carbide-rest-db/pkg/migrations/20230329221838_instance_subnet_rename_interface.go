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

func renameInstanceSubnetToInterfaceUpMigration(ctx context.Context, db *bun.DB) error {
	// Start transaction
	tx, terr := db.BeginTx(ctx, &sql.TxOptions{})
	if terr != nil {
		handlePanic(terr, "failed to begin transaction")
		return terr
	}

	// Rename instance_subnet table to interface if it exists
	_, err := tx.Exec("ALTER TABLE IF EXISTS instance_subnet RENAME TO interface")
	handleError(tx, err)

	// Drop the older index if exists
	_, err = tx.Exec("DROP INDEX IF EXISTS instance_subnet_status_idx")
	handleError(tx, err)

	// Drop if the index exists (won't occur/harmless in dev/stage/prod but helps with test)
	_, err = tx.Exec("DROP INDEX IF EXISTS interface_status_idx")
	handleError(tx, err)

	// Add status index for interface (formerly instance_subnet) model
	_, err = tx.Exec("CREATE INDEX interface_status_idx ON public.interface(status) WHERE deleted IS NULL")
	handleError(tx, err)

	terr = tx.Commit()
	if terr != nil {
		handlePanic(terr, "failed to commit transaction")
		return terr
	}

	fmt.Print(" [up migration] ")
	return nil
}

func init() {
	Migrations.MustRegister(renameInstanceSubnetToInterfaceUpMigration, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
