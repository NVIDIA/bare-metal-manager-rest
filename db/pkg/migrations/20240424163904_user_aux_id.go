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
			fmt.Println("failed to begin transaction: ", terr)
			return terr
		}

		// Add AuxID column to User table
		_, err := tx.NewAddColumn().Model((*model.User)(nil)).IfNotExists().ColumnExpr("auxiliary_id TEXT").Exec(ctx)
		if err != nil {
			fmt.Println("failed to add column: ", err)
			return err
		}

		// Drop if the index exists (won't occur/harmless in dev/stage/prod but helps with test)
		_, err = tx.Exec("DROP INDEX IF EXISTS user_auxiliary_id_indx")
		handleError(tx, err)

		_, err = tx.Exec("CREATE INDEX user_auxiliary_id_indx ON public.user(auxiliary_id)")
		handleError(tx, err)

		terr = tx.Commit()
		if terr != nil {
			handlePanic(terr, "failed to commit transaction")
		}
		fmt.Print(" [up migration] Added 'auxiliary_id' column to 'user' table successfully and created index 'user_auxiliary_id_indx'")

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		// Start transactions
		tx, terr := db.BeginTx(ctx, &sql.TxOptions{})
		if terr != nil {
			fmt.Println("failed to begin transaction: ", terr)
			return terr
		}

		// Remove AuxID column from User table
		_, err := tx.NewDropColumn().Model((*model.User)(nil)).Column("auxiliary_id").Exec(ctx)
		if err != nil {
			fmt.Println("failed to drop column: ", err)
			return err
		}

		terr = tx.Commit()
		if terr != nil {
			fmt.Println("failed to commit transaction: ", terr)
			return terr
		}

		fmt.Print(" [down migration] Removed 'aux_id' column from 'user' table successfully.")
		return nil
	})
}
