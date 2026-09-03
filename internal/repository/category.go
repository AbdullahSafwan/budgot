package repository

import (
	"context"

	"budgot/internal/configs"
	"budgot/internal/ent"
	"budgot/internal/ent/category"
	"budgot/internal/ent/predicate"
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

func (r *CategoryRepository) FindByID(ctx context.Context, ownerID, id int) (*ent.Category, error) {
	c, err := r.client.Category.Query().
		Where(category.IDEQ(id), category.HasOwnerWith(user.IDEQ(ownerID))).
		Only(ctx)
	return c, configs.Translate(err)
}

type ListCategoriesParams struct {
	OwnerID       int
	Type          *category.Type
	Limit, Offset int
}

// List is not country-scoped: categories are deliberately shared across countries.
func (r *CategoryRepository) List(ctx context.Context, p ListCategoriesParams) ([]*ent.Category, error) {
	preds := []predicate.Category{
		category.HasOwnerWith(user.IDEQ(p.OwnerID)),
		category.IsActiveEQ(true),
	}
	if p.Type != nil {
		preds = append(preds, category.TypeEQ(*p.Type))
	}

	return r.client.Category.Query().
		Where(preds...).
		Order(ent.Asc(category.FieldName)).
		Limit(listLimit(p.Limit)).
		Offset(p.Offset).
		All(ctx)
}

// Delete soft-deletes a category; the name becomes reusable once inactive.
func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	_, err := r.client.Category.UpdateOneID(id).SetIsActive(false).Save(ctx)
	return configs.Translate(err)
}
