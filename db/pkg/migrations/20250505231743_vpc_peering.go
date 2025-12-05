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

		_, err := tx.NewCreateTable().Model((*model.VpcPeering)(nil)).Exec(ctx)
		handleError(tx, err)

		// Create VpcPeering index on all vpc1_id, vpc2_id, and site_id
		_, err = tx.Exec("DROP INDEX IF EXISTS idx_vpc_peering_vpc1")
		handleError(tx, err)
		_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_vpc_peering_vpc1 ON vpc_peering(vpc1_id)")
		handleError(tx, err)

		_, err = tx.Exec("DROP INDEX IF EXISTS idx_vpc_peering_vpc2")
		handleError(tx, err)
		_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_vpc_peering_vpc2 ON vpc_peering(vpc2_id)")
		handleError(tx, err)

		_, err = tx.Exec("DROP INDEX IF EXISTS idx_vpc_peering_site_id")
		handleError(tx, err)
		_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_vpc_peering_site_id ON vpc_peering(site_id)")
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
