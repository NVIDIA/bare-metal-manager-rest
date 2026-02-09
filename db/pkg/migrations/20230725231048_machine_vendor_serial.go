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

		// Add vendor column to Machine table
		_, err := tx.NewAddColumn().Model((*model.Machine)(nil)).IfNotExists().ColumnExpr("vendor VARCHAR").Exec(ctx)
		handleError(tx, err)

		// Add product_name column to Machine table
		_, err = tx.NewAddColumn().Model((*model.Machine)(nil)).IfNotExists().ColumnExpr("product_name VARCHAR").Exec(ctx)
		handleError(tx, err)

		// Add serial_number column to Machine table
		_, err = tx.NewAddColumn().Model((*model.Machine)(nil)).IfNotExists().ColumnExpr("serial_number VARCHAR").Exec(ctx)
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
