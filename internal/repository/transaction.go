package repository

import (
	"context"
	"time"

	"budgot/internal/ent"
	"budgot/internal/ent/account"
	"budgot/internal/ent/category"
	"budgot/internal/ent/country"
	"budgot/internal/ent/currency"
	"budgot/internal/ent/transaction"
)

type TransactionRepository struct {
	client *ent.Client
}

func NewTransactionRepository(client *ent.Client) *TransactionRepository {
	return &TransactionRepository{client: client}
}

type CreateTransactionParams struct {
	OwnerID         int
	AccountID       int
	CategoryID      int
	Amount          int64
	Description     string
	TransactionDate time.Time
}

func (r *TransactionRepository) Create(ctx context.Context, params CreateTransactionParams) (*ent.Transaction, error) {
	return r.client.Transaction.Create().
		SetOwnerID(params.OwnerID).
		SetAccountID(params.AccountID).
		SetCategoryID(params.CategoryID).
		SetAmount(params.Amount).
		SetDescription(params.Description).
		SetTransactionDate(params.TransactionDate).
		Save(ctx)
}

func (r *TransactionRepository) FindByID(ctx context.Context, id int) (*ent.Transaction, error) {
	return r.client.Transaction.Query().Where(transaction.IDEQ(id)).Only(ctx)
}

func (r *TransactionRepository) ListByAccount(ctx context.Context, accountID int) ([]*ent.Transaction, error) {
	return r.client.Transaction.Query().
		Where(transaction.HasAccountWith(account.IDEQ(accountID))).
		All(ctx)
}

// links a transaction to another transaction (e.g. for transfers)
func (r *TransactionRepository) SetLinkedTransaction(ctx context.Context, id, linkedID int) error {
	_, err := r.client.Transaction.UpdateOneID(id).SetLinkedTransactionID(linkedID).Save(ctx)
	return err
}

func (r *TransactionRepository) SumByAccount(ctx context.Context, accountID int) (int64, error) {
	sum, err := r.client.Transaction.Query().
		Where(transaction.HasAccountWith(account.IDEQ(accountID))).
		Aggregate(ent.Sum(transaction.FieldAmount)).
		Int(ctx)
	if err != nil {
		return 0, err
	}
	return int64(sum), nil
}

// sums trx for a category per month/year, filtered by country and currency of the account
func (r *TransactionRepository) SumByCategoryAndPeriod(ctx context.Context, categoryID, countryID, currencyID, month, year int) (int64, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	sum, err := r.client.Transaction.Query().
		Where(
			transaction.HasCategoryWith(category.IDEQ(categoryID)),
			transaction.HasAccountWith(
				account.HasCountryWith(country.IDEQ(countryID)),
				account.HasCurrencyWith(currency.IDEQ(currencyID)),
			),
			transaction.TransactionDateGTE(start),
			transaction.TransactionDateLT(end),
		).
		Aggregate(ent.Sum(transaction.FieldAmount)).
		Int(ctx)
	if err != nil {
		return 0, err
	}
	return int64(sum), nil
}
