package configs

import (
	"errors"

	"budgot/internal/ent"
)

// Domain errors, translated from raw Ent errors.
var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

// Translate maps a raw Ent error to a domain error.
func Translate(err error) error {
	switch {
	case err == nil:
		return nil
	case ent.IsNotFound(err):
		return ErrNotFound
	case ent.IsConstraintError(err):
		return ErrConflict
	default:
		return err
	}
}
