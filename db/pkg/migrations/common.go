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
