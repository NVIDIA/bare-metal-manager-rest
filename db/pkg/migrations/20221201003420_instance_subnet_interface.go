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

		// Add is_physical column to Interface (formerly InstanceSubnet) model
		_, err := tx.NewAddColumn().Model((*model.Interface)(nil)).IfNotExists().ColumnExpr("is_physical BOOLEAN NOT NULL DEFAULT true").Exec(ctx)
		handleError(tx, err)

		// Add virtual_function_id column to Interface (formerly InstanceSubnet) model
		_, err = tx.NewAddColumn().Model((*model.Interface)(nil)).IfNotExists().ColumnExpr("virtual_function_id INTEGER").Exec(ctx)
		handleError(tx, err)

		// Add mac_address column to Interface (formerly InstanceSubnet) model
		_, err = tx.NewAddColumn().Model((*model.Interface)(nil)).IfNotExists().ColumnExpr("mac_address VARCHAR").Exec(ctx)
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
