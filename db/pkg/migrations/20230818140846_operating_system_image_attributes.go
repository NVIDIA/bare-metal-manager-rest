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

func operatingSystemImageAttributeUpMigration(ctx context.Context, db *bun.DB) error {
	// Start transactions
	tx, terr := db.BeginTx(ctx, &sql.TxOptions{})
	if terr != nil {
		handlePanic(terr, "failed to begin transaction")
	}

	// Add type column to OperatingSystem table
	_, err := tx.NewAddColumn().Model((*model.OperatingSystem)(nil)).IfNotExists().ColumnExpr("type varchar").Exec(ctx)
	handleError(tx, err)

	// Add image attributes columns to OperatingSystem table
	_, err = tx.NewAddColumn().Model((*model.OperatingSystem)(nil)).IfNotExists().ColumnExpr("image_sha varchar").Exec(ctx)
	handleError(tx, err)

	_, err = tx.NewAddColumn().Model((*model.OperatingSystem)(nil)).IfNotExists().ColumnExpr("image_auth_type varchar").Exec(ctx)
	handleError(tx, err)

	_, err = tx.NewAddColumn().Model((*model.OperatingSystem)(nil)).IfNotExists().ColumnExpr("image_auth_token varchar").Exec(ctx)
	handleError(tx, err)

	_, err = tx.NewAddColumn().Model((*model.OperatingSystem)(nil)).IfNotExists().ColumnExpr("image_disk varchar").Exec(ctx)
	handleError(tx, err)

	_, err = tx.NewAddColumn().Model((*model.OperatingSystem)(nil)).IfNotExists().ColumnExpr("root_fs_id varchar").Exec(ctx)
	handleError(tx, err)

	// Update Type record for each Operating System
	res, err := tx.Exec("UPDATE operating_system SET type = 'iPXE'")
	handleError(tx, err)

	osRowAffected, _ := res.RowsAffected()
	fmt.Printf("Updated %v operating systems \n", osRowAffected)

	// Make type column not nullable after we updated each row to default value
	_, err = tx.Exec("ALTER TABLE operating_system ALTER COLUMN type SET NOT NULL")
	handleError(tx, err)

	terr = tx.Commit()
	if terr != nil {
		handlePanic(terr, "failed to commit transaction")
	}

	fmt.Print(" [up migration] ")
	return nil
}

func init() {
	Migrations.MustRegister(operatingSystemImageAttributeUpMigration, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
