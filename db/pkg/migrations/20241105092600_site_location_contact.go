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

		// Add location column
		_, err := tx.NewAddColumn().Model((*model.Site)(nil)).IfNotExists().ColumnExpr("location JSONB").Exec(ctx)
		handleError(tx, err)

		// Add contact column
		_, err = tx.NewAddColumn().Model((*model.Site)(nil)).IfNotExists().ColumnExpr("contact JSONB").Exec(ctx)
		handleError(tx, err)

		// update tsv index
		_, err = tx.Exec("DROP INDEX IF EXISTS site_tsv_idx")
		handleError(tx, err)

		// Add tsv index with new columns
		_, err = tx.Exec("CREATE INDEX site_tsv_idx ON site USING gin(to_tsvector('english', name || ' ' || description || ' ' || status || ' ' || location::text || ' ' || contact::text))")
		handleError(tx, err)

		terr = tx.Commit()
		if terr != nil {
			handlePanic(terr, "failed to commit transaction")
		}

		fmt.Print(" [up migration] Added 'location' and 'contact' columns to 'site' table successfully. ")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
