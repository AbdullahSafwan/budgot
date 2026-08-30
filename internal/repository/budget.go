package repository

import (
	"context"

	"budgot/internal/ent"
	"budgot/internal/ent/budget"
	"budgot/internal/ent/user"
)

type BudgetRepository struct {
	client *ent.Client
}

func NewBudgetRepository(client *ent.Client) *BudgetRepository {
	return &BudgetRepository{client: client}
}

// WithTx returns a copy of the repository bound to the given transaction, so its
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
	return r.client.Budget.Create().
		SetOwnerID(params.OwnerID).
		SetCategoryID(params.CategoryID).
		SetCountryID(params.CountryID).
		SetCurrencyID(params.CurrencyID).
		SetMonth(params.Month).
		SetYear(params.Year).
		SetAmount(params.Amount).
		Save(ctx)
}

func (r *BudgetRepository) FindByID(ctx context.Context, id int) (*ent.Budget, error) {
	return r.client.Budget.Query().Where(budget.IDEQ(id)).Only(ctx)
}

func (r *BudgetRepository) ListByOwner(ctx context.Context, ownerID int) ([]*ent.Budget, error) {
	return r.client.Budget.Query().
		Where(budget.HasOwnerWith(user.IDEQ(ownerID))).
		All(ctx)
}
