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

		// Add agent_cert_expiry column to the Site table
		_, err := tx.NewAddColumn().Model((*model.Site)(nil)).IfNotExists().ColumnExpr("agent_cert_expiry TIMESTAMPTZ DEFAULT NULL").Exec(ctx)
		handleError(tx, err)

		terr = tx.Commit()
		if terr != nil {
			handlePanic(terr, "failed to commit transaction")
		}

		fmt.Print("  [up migration] Added agent_cert_expiry column to Site table")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		// Rollback script to remove cert_expiration column
		_, err := db.NewDropColumn().Model((*model.Site)(nil)).Column("agent_cert_expiry").Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to drop agent_cert_expiry column: %v", err)
		}

		fmt.Print(" [down migration] Removed agent_cert_expiry column from Site table\n")
		return nil
	})
}
