package repository

import (
	"context"

	"budgot/internal/ent"
	"budgot/internal/ent/category"
	"budgot/internal/ent/user"
)

type CategoryRepository struct {
	client *ent.Client
}

func NewCategoryRepository(client *ent.Client) *CategoryRepository {
	return &CategoryRepository{client: client}
}

type CreateCategoryParams struct {
	OwnerID int
	Name    string
	Type    category.Type
	Color   string
}

func (r *CategoryRepository) Create(ctx context.Context, params CreateCategoryParams) (*ent.Category, error) {
	return r.client.Category.Create().
		SetOwnerID(params.OwnerID).
		SetName(params.Name).
		SetType(params.Type).
		SetColor(params.Color).
		Save(ctx)
}

func (r *CategoryRepository) FindByID(ctx context.Context, id int) (*ent.Category, error) {
	return r.client.Category.Query().Where(category.IDEQ(id)).Only(ctx)
}

func (r *CategoryRepository) ListByOwner(ctx context.Context, ownerID int) ([]*ent.Category, error) {
	return r.client.Category.Query().
		Where(category.HasOwnerWith(user.IDEQ(ownerID))).
		All(ctx)
}
