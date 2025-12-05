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
	"fmt"

	"github.com/uptrace/bun"
)

func handleError(tx bun.Tx, err error) {
	if err == nil {
		return
	}

	terr := tx.Rollback()
	if terr != nil {
		handlePanic(terr, "failed to rollback transaction")
	}

	handlePanic(err, "failed to execute migration")
}

func handlePanic(err error, message string) {
	if err != nil {
		fmt.Printf("unrecoverable error: %v, details: %v", message, err)
		panic(err)
	}
}
