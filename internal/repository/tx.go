package repository

import (
	"context"
	"errors"

	"budgot/internal/ent"
)

// defaultListLimit caps unbounded list queries when the caller doesn't set Limit.
const defaultListLimit = 100

// listLimit returns limit if positive, else defaultListLimit.
func listLimit(limit int) int {
	if limit <= 0 {
		return defaultListLimit
	}
	return limit
}

// WithTx runs fn in a transaction, committing on success and rolling back on error.
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
