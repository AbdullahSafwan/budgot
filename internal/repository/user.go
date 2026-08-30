package repository

import (
	"context"

	"budgot/internal/ent"
	"budgot/internal/ent/user"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

// WithTx returns a copy of the repository bound to the given transaction, so its
func (r *UserRepository) WithTx(tx *ent.Tx) *UserRepository {
	return &UserRepository{client: tx.Client()}
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*ent.User, error) {
	return r.client.User.Query().Where(user.UsernameEQ(username)).Only(ctx)
}
