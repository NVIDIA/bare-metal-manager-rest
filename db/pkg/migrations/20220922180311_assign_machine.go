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

		// Add "is_assigned" column to Machine table
		_, err := tx.NewAddColumn().Model((*model.Machine)(nil)).IfNotExists().ColumnExpr("is_assigned BOOLEAN NOT NULL DEFAULT false").Exec(ctx)
		handleError(tx, err)

		// Add "allocation_constraint_id" column to Instance table
		_, err = tx.NewAddColumn().Model((*model.Instance)(nil)).IfNotExists().ColumnExpr("allocation_constraint_id uuid NOT NULL").Exec(ctx)
		handleError(tx, err)

		// Add allocation contraint foreign key
		// Drop if one exists (won't occur/harmless in dev/stage/prod but helps with test)
		_, err = tx.Exec("ALTER TABLE instance DROP CONSTRAINT IF EXISTS instance_allocation_constraint_id_fkey")
		handleError(tx, err)

		_, err = tx.Exec("ALTER TABLE instance ADD CONSTRAINT instance_allocation_constraint_id_fkey FOREIGN KEY (allocation_constraint_id) REFERENCES public.allocation_constraint(id)")
		handleError(tx, err)

		terr = tx.Commit()
		if terr != nil {
			handlePanic(terr, "failed to commit transaction")
		}

		fmt.Print(" [up migration] ")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
