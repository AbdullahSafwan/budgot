package repository

import (
	"context"

	"budgot/internal/configs"
	"budgot/internal/ent"
	"budgot/internal/ent/user"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

// WithTx binds the repository to an existing transaction.
func (r *UserRepository) WithTx(tx *ent.Tx) *UserRepository {
	return &UserRepository{client: tx.Client()}
}

// FindByUsername only returns active users, so a deactivated account fails login like a bad username.
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*ent.User, error) {
	u, err := r.client.User.Query().
		Where(user.UsernameEQ(username), user.IsActiveEQ(true)).
		Only(ctx)
	return u, configs.Translate(err)
}

// Delete soft-deletes a user; rows are never hard-deleted.
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	_, err := r.client.User.UpdateOneID(id).SetIsActive(false).Save(ctx)
	return configs.Translate(err)
}
