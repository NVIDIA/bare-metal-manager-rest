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

		// Add is_infinity_enabled column to site table
		_, err := tx.NewAddColumn().Model((*model.Site)(nil)).IfNotExists().ColumnExpr("is_infinity_enabled BOOLEAN NOT NULL DEFAULT false").Exec(ctx)
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
