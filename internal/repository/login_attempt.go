package repository

import (
	"context"

	"budgot/internal/configs"
	"budgot/internal/ent"
)

type LoginAttemptRepository struct {
	client *ent.Client
}

func NewLoginAttemptRepository(client *ent.Client) *LoginAttemptRepository {
	return &LoginAttemptRepository{client: client}
}

// WithTx binds the repository to an existing transaction.
func (r *LoginAttemptRepository) WithTx(tx *ent.Tx) *LoginAttemptRepository {
	return &LoginAttemptRepository{client: tx.Client()}
}

func (r *LoginAttemptRepository) Record(ctx context.Context, username, ipAddress string, success bool) error {
	_, err := r.client.LoginAttempt.Create().
		SetUsername(username).
		SetIPAddress(ipAddress).
		SetSuccess(success).
		Save(ctx)
	return configs.Translate(err)
}
