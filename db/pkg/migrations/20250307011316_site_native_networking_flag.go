package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/nvidia/carbide-rest/db/pkg/db/model"
)

// TODO: Add migration to remove the column if it exists
func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		// Start transactions
		tx, terr := db.BeginTx(ctx, &sql.TxOptions{})
		if terr != nil {
			handlePanic(terr, "failed to begin transaction")
		}

		// Add is_native_networking_enabled column
		_, err := tx.NewAddColumn().Model((*model.Site)(nil)).IfNotExists().ColumnExpr("is_native_networking_enabled bool DEFAULT false").Exec(ctx)
		handleError(tx, err)

		terr = tx.Commit()
		if terr != nil {
			handlePanic(terr, "failed to commit transaction")
		}

		fmt.Print(" [up migration] Added 'is_native_networking_enabled' to 'site' table successfully. ")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
