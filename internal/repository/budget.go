package repository

import (
	"context"

	"budgot/internal/configs"
	"budgot/internal/ent"
	"budgot/internal/ent/budget"
	"budgot/internal/ent/country"
	"budgot/internal/ent/user"
)

type BudgetRepository struct {
	client *ent.Client
}

func NewBudgetRepository(client *ent.Client) *BudgetRepository {
	return &BudgetRepository{client: client}
}

// WithTx binds the repository to an existing transaction.
func (r *BudgetRepository) WithTx(tx *ent.Tx) *BudgetRepository {
	return &BudgetRepository{client: tx.Client()}
}

type CreateBudgetParams struct {
	OwnerID    int
	CategoryID int
	CountryID  int
	CurrencyID int
	Month      int
	Year       int
	Amount     int64
}

func (r *BudgetRepository) Create(ctx context.Context, params CreateBudgetParams) (*ent.Budget, error) {
	b, err := r.client.Budget.Create().
		SetOwnerID(params.OwnerID).
		SetCategoryID(params.CategoryID).
		SetCountryID(params.CountryID).
		SetCurrencyID(params.CurrencyID).
		SetMonth(params.Month).
		SetYear(params.Year).
		SetAmount(params.Amount).
		Save(ctx)
	return b, configs.Translate(err)
}

func (r *BudgetRepository) FindByID(ctx context.Context, id int) (*ent.Budget, error) {
	b, err := r.client.Budget.Query().Where(budget.IDEQ(id)).Only(ctx)
	return b, configs.Translate(err)
}

// list function for filtering budgets by owner and country
func (r *BudgetRepository) ListByOwnerAndCountry(ctx context.Context, ownerID, countryID int) ([]*ent.Budget, error) {
	return r.client.Budget.Query().
		Where(
			budget.HasOwnerWith(user.IDEQ(ownerID)),
			budget.HasCountryWith(country.IDEQ(countryID)),
			budget.IsActiveEQ(true),
		).
		All(ctx)
}

// Delete soft-deletes a budget; rows are never hard-deleted.
func (r *BudgetRepository) Delete(ctx context.Context, id int) error {
	_, err := r.client.Budget.UpdateOneID(id).SetIsActive(false).Save(ctx)
	return configs.Translate(err)
}
