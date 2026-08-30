package repository

import (
	"context"
	"errors"

	"budgot/internal/ent"
)

// atomically executes the given function within a transaction. If the function returns an error, the transaction is rolled back; otherwise, it is committed.
func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return errors.Join(err, rerr)
		}
		return err
	}
	return tx.Commit()
}
