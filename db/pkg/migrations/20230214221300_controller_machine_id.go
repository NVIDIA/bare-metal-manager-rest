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

func machineControllerMachineIDUpMigration(ctx context.Context, db *bun.DB) error {
	// Start transactions
	tx, terr := db.BeginTx(ctx, &sql.TxOptions{})
	if terr != nil {
		handlePanic(terr, "failed to begin transaction")
	}

	// Alter the type of controller_machine_id field in the machine table from uuid to string
	_, err := tx.Exec("ALTER TABLE machine ALTER COLUMN controller_machine_id TYPE TEXT USING controller_machine_id::text")
	handleError(tx, err)

	terr = tx.Commit()
	if terr != nil {
		handlePanic(terr, "failed to commit transaction")
	}

	fmt.Print(" [up migration] ")
	return nil
}

func init() {
	Migrations.MustRegister(machineControllerMachineIDUpMigration, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
