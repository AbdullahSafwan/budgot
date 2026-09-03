package repository

import (
	"context"

	"budgot/internal/configs"
	"budgot/internal/ent"
	"budgot/internal/ent/budget"
	"budgot/internal/ent/category"
	"budgot/internal/ent/country"
	"budgot/internal/ent/predicate"
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

func (r *BudgetRepository) FindByID(ctx context.Context, ownerID, id int) (*ent.Budget, error) {
	b, err := r.client.Budget.Query().
		Where(budget.IDEQ(id), budget.HasOwnerWith(user.IDEQ(ownerID))).
		Only(ctx)
	return b, configs.Translate(err)
}

type ListBudgetsParams struct {
	OwnerID       int
	CountryID     int
	CategoryID    *int
	Month, Year   *int
	WithEdges     bool
	Limit, Offset int
}

func (r *BudgetRepository) List(ctx context.Context, p ListBudgetsParams) ([]*ent.Budget, error) {
	preds := []predicate.Budget{
		budget.HasOwnerWith(user.IDEQ(p.OwnerID)),
		budget.HasCountryWith(country.IDEQ(p.CountryID)),
		budget.IsActiveEQ(true),
	}
	if p.CategoryID != nil {
		preds = append(preds, budget.HasCategoryWith(category.IDEQ(*p.CategoryID)))
	}
	if p.Month != nil {
		preds = append(preds, budget.MonthEQ(*p.Month))
	}
	if p.Year != nil {
		preds = append(preds, budget.YearEQ(*p.Year))
	}

	q := r.client.Budget.Query().
		Where(preds...).
		Limit(listLimit(p.Limit)).
		Offset(p.Offset)
	if p.WithEdges {
		q = q.WithCategory().WithCurrency()
	}
	return q.All(ctx)
}

// Delete soft-deletes a budget; rows are never hard-deleted.
func (r *BudgetRepository) Delete(ctx context.Context, id int) error {
	_, err := r.client.Budget.UpdateOneID(id).SetIsActive(false).Save(ctx)
	return configs.Translate(err)
}
