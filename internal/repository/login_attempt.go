package repository

import (
	"context"

	"budgot/internal/ent"
)

type LoginAttemptRepository struct {
	client *ent.Client
}

func NewLoginAttemptRepository(client *ent.Client) *LoginAttemptRepository {
	return &LoginAttemptRepository{client: client}
}

// WithTx returns a copy of the repository bound to the given transaction, so its
func (r *LoginAttemptRepository) WithTx(tx *ent.Tx) *LoginAttemptRepository {
	return &LoginAttemptRepository{client: tx.Client()}
}

func (r *LoginAttemptRepository) Record(ctx context.Context, username, ipAddress string, success bool) error {
	_, err := r.client.LoginAttempt.Create().
		SetUsername(username).
		SetIPAddress(ipAddress).
		SetSuccess(success).
		Save(ctx)
	return err
}
