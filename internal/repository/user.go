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

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*ent.User, error) {
	return r.client.User.Query().Where(user.UsernameEQ(username)).Only(ctx)
}
