package repository

import (
	"context"

	"budgot/internal/configs"
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

// WithTx binds the repository to an existing transaction.
func (r *CategoryRepository) WithTx(tx *ent.Tx) *CategoryRepository {
	return &CategoryRepository{client: tx.Client()}
}

type CreateCategoryParams struct {
	OwnerID int
	Name    string
	Type    category.Type
	Color   string
}

func (r *CategoryRepository) Create(ctx context.Context, params CreateCategoryParams) (*ent.Category, error) {
	c, err := r.client.Category.Create().
		SetOwnerID(params.OwnerID).
		SetName(params.Name).
		SetType(params.Type).
		SetColor(params.Color).
		Save(ctx)
	return c, configs.Translate(err)
}

func (r *CategoryRepository) FindByID(ctx context.Context, id int) (*ent.Category, error) {
	c, err := r.client.Category.Query().Where(category.IDEQ(id)).Only(ctx)
	return c, configs.Translate(err)
}

func (r *CategoryRepository) ListByOwner(ctx context.Context, ownerID int) ([]*ent.Category, error) {
	return r.client.Category.Query().
		Where(
			category.HasOwnerWith(user.IDEQ(ownerID)),
			category.IsActiveEQ(true),
		).
		All(ctx)
}

// Delete soft-deletes a category; the name becomes reusable once inactive.
func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	_, err := r.client.Category.UpdateOneID(id).SetIsActive(false).Save(ctx)
	return configs.Translate(err)
}
