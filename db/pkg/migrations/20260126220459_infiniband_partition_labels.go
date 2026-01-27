/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
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

	"github.com/nvidia/carbide-rest/db/pkg/db/model"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		// Start transactions
		tx, terr := db.BeginTx(ctx, &sql.TxOptions{})
		if terr != nil {
			handlePanic(terr, "failed to begin transaction")
		}

		// Add labels column to InfiniBandPartition table
		_, err := tx.NewAddColumn().Model((*model.InfiniBandPartition)(nil)).IfNotExists().ColumnExpr("labels JSONB NOT NULL DEFAULT ('{}')").Exec(ctx)
		handleError(tx, err)

		// Drop if the existing infiniband_partition_gin_idx exists
		_, err = tx.Exec("DROP INDEX IF EXISTS infiniband_partition_gin_idx")
		handleError(tx, err)

		// Recreate the GIN index with labels included for text search
		_, err = tx.Exec("CREATE INDEX infiniband_partition_gin_idx ON public.infiniband_partition USING GIN (name gin_trgm_ops, description gin_trgm_ops, partition_key gin_trgm_ops, partition_name gin_trgm_ops, status gin_trgm_ops, labels gin_trgm_ops)")
		handleError(tx, err)

		terr = tx.Commit()
		if terr != nil {
			handlePanic(terr, "failed to commit transaction")
		}

		fmt.Print(" [up migration] Added 'labels' column to 'infiniband_partition' table successfully. ")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
