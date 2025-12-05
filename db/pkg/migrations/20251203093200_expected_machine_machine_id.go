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
	"github.com/nvidia/carbide-rest/db/pkg/db/model"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		// Start transactions
		tx, terr := db.BeginTx(ctx, &sql.TxOptions{})
		if terr != nil {
			handlePanic(terr, "failed to begin transaction")
		}

		// Add machine_id column to expected_machine table
		_, err := tx.NewAddColumn().Model((*model.ExpectedMachine)(nil)).IfNotExists().ColumnExpr("machine_id TEXT NULL").Exec(ctx)
		handleError(tx, err)

		// Add Machine foreign key for expected_machine
		// Drop if one exists (won't occur/harmless in dev/stage/prod but helps with test)
		_, err = tx.Exec("ALTER TABLE expected_machine DROP CONSTRAINT IF EXISTS expected_machine_machine_id_fkey")
		handleError(tx, err)

		_, err = tx.Exec("ALTER TABLE expected_machine ADD CONSTRAINT expected_machine_machine_id_fkey FOREIGN KEY (machine_id) REFERENCES public.machine(id)")
		handleError(tx, err)

		// Drop index if it exists
		_, err = tx.Exec("DROP INDEX IF EXISTS expected_machine_machine_id_idx")
		handleError(tx, err)

		// Add index for machine_id for improved query performance
		_, err = tx.Exec("CREATE INDEX expected_machine_machine_id_idx ON expected_machine(machine_id)")
		handleError(tx, err)

		terr = tx.Commit()
		if terr != nil {
			handlePanic(terr, "failed to commit transaction")
		}

		fmt.Print(" [up migration] Added machine_id column and foreign key to expected_machine table. ")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] No action taken")
		return nil
	})
}

