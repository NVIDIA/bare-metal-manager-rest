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
		tx, terr := db.BeginTx(ctx, &sql.TxOptions{})
		if terr != nil {
			handlePanic(terr, "failed to begin transaction")
		}

		// create SKU table
		_, err := tx.NewCreateTable().Model((*model.SKU)(nil)).IfNotExists().Exec(ctx)
		handleError(tx, err)

		// index for sku_id
		_, err = tx.Exec("DROP INDEX IF EXISTS sku_id_idx")
		handleError(tx, err)
		_, err = tx.Exec("CREATE INDEX sku_id_idx ON sku(id)")
		handleError(tx, err)

		// index for site_id
		_, err = tx.Exec("DROP INDEX IF EXISTS sku_site_id_idx")
		handleError(tx, err)
		_, err = tx.Exec("CREATE INDEX sku_site_id_idx ON sku(site_id)")
		handleError(tx, err)

		// index for created
		_, err = tx.Exec("DROP INDEX IF EXISTS sku_created_idx")
		handleError(tx, err)
		_, err = tx.Exec("CREATE INDEX sku_created_idx ON sku(created)")
		handleError(tx, err)

		// index for updated
		_, err = tx.Exec("DROP INDEX IF EXISTS sku_updated_idx")
		handleError(tx, err)
		_, err = tx.Exec("CREATE INDEX sku_updated_idx ON sku(updated)")
		handleError(tx, err)

		terr = tx.Commit()
		if terr != nil {
			handlePanic(terr, "failed to commit transaction")
		}

		fmt.Print(" [up migration] ")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] No action taken")
		return nil
	})
}
