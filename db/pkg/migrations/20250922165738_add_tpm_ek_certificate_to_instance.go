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

		// Add tpm_ek_certificate column to instance table
		_, err := tx.NewAddColumn().Model((*model.Instance)(nil)).IfNotExists().ColumnExpr("tpm_ek_certificate TEXT").Exec(ctx)
		handleError(tx, err)

		terr = tx.Commit()
		if terr != nil {
			handlePanic(terr, "failed to commit transaction")
		}

		fmt.Print(" [up migration] Added 'tpm_ek_certificate' to 'instance' table successfully. ")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
