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

func siteSshHostnameRenameUpMigration(ctx context.Context, db *bun.DB) error {
	// Start transaction
	tx, terr := db.BeginTx(ctx, &sql.TxOptions{})
	if terr != nil {
		handlePanic(terr, "failed to begin transaction")
	}

	// Rename site sshHostname column to serialConsoleHostname if that column exists
	res, err := tx.Exec("SELECT column_name FROM information_schema.columns WHERE table_name = 'site' AND column_name = 'ssh_hostname'")
	handleError(tx, err)
	sshHostnameRowsAffected, err := res.RowsAffected()
	res, err = tx.Exec("SELECT column_name FROM information_schema.columns WHERE table_name = 'site' AND column_name = 'serial_console_hostname'")
	handleError(tx, err)
	serialConsoleHostnameRowsAffected, err := res.RowsAffected()
	handleError(tx, err)
	if sshHostnameRowsAffected > 0 && serialConsoleHostnameRowsAffected == 0 {
		_, err := tx.Exec("ALTER TABLE site RENAME COLUMN ssh_hostname TO serial_console_hostname")
		handleError(tx, err)
	} else {
		fmt.Println("sshHostname rename to serialConsoleHostname: Migration skipped. Either the column does not exist or already renamed")
	}

	terr = tx.Commit()
	if terr != nil {
		handlePanic(terr, "failed to commit transaction")
	}

	fmt.Print(" [up migration] ")
	return nil
}

func init() {
	Migrations.MustRegister(siteSshHostnameRenameUpMigration, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
