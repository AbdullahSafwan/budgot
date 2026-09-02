package repository

import (
	"context"
	"time"

	"budgot/internal/configs"
	"budgot/internal/ent"
	"budgot/internal/ent/account"
	"budgot/internal/ent/category"
	"budgot/internal/ent/country"
	"budgot/internal/ent/currency"
	"budgot/internal/ent/transaction"
	"budgot/internal/ent/user"
)

type TransactionRepository struct {
	client *ent.Client
}

func NewTransactionRepository(client *ent.Client) *TransactionRepository {
	return &TransactionRepository{client: client}
}

// WithTx binds the repository to an existing transaction.
func (r *TransactionRepository) WithTx(tx *ent.Tx) *TransactionRepository {
	return &TransactionRepository{client: tx.Client()}
}

type CreateTransactionParams struct {
	OwnerID         int
	AccountID       int
	CategoryID      int
	CountryID       int
	Amount          int64
	Description     string
	TransactionDate time.Time
	TransferGroup   *string
}

func (r *TransactionRepository) Create(ctx context.Context, params CreateTransactionParams) (*ent.Transaction, error) {
	c := r.client.Transaction.Create().
		SetOwnerID(params.OwnerID).
		SetAccountID(params.AccountID).
		SetCategoryID(params.CategoryID).
		SetCountryID(params.CountryID).
		SetAmount(params.Amount).
		SetDescription(params.Description).
		SetTransactionDate(params.TransactionDate)
	if params.TransferGroup != nil {
		c = c.SetTransferGroup(*params.TransferGroup)
	}
	t, err := c.Save(ctx)
	return t, configs.Translate(err)
}

func (r *TransactionRepository) FindByID(ctx context.Context, id int) (*ent.Transaction, error) {
	t, err := r.client.Transaction.Query().Where(transaction.IDEQ(id)).Only(ctx)
	return t, configs.Translate(err)
}

func (r *TransactionRepository) ListByAccount(ctx context.Context, accountID int) ([]*ent.Transaction, error) {
	return r.client.Transaction.Query().
		Where(transaction.HasAccountWith(account.IDEQ(accountID))).
		All(ctx)
}

func (r *TransactionRepository) ListByOwnerAndCountry(ctx context.Context, ownerID, countryID int) ([]*ent.Transaction, error) {
	return r.client.Transaction.Query().
		Where(
			transaction.HasOwnerWith(user.IDEQ(ownerID)),
			transaction.HasCountryWith(country.IDEQ(countryID)),
		).
		All(ctx)
}

// SetTransferGroup links a transaction to another as one leg of a transfer.
func (r *TransactionRepository) SetTransferGroup(ctx context.Context, id int, group string) error {
	_, err := r.client.Transaction.UpdateOneID(id).SetTransferGroup(group).Save(ctx)
	return configs.Translate(err)
}

func (r *TransactionRepository) ListByTransferGroup(ctx context.Context, group string) ([]*ent.Transaction, error) {
	return r.client.Transaction.Query().
		Where(transaction.TransferGroupEQ(group)).
		All(ctx)
}

func (r *TransactionRepository) SumByAccount(ctx context.Context, accountID int) (int64, error) {
	var result []struct {
		Total *int64 `json:"total"`
	}
	err := r.client.Transaction.Query().
		Where(transaction.HasAccountWith(account.IDEQ(accountID))).
		Aggregate(ent.As(ent.Sum(transaction.FieldAmount), "total")).
		Scan(ctx, &result)
	if err != nil {
		return 0, err
	}
	if len(result) == 0 || result[0].Total == nil {
		return 0, nil
	}
	return *result[0].Total, nil
}

// sums trx for a category per month/year, filtered by country and currency of the account
func (r *TransactionRepository) SumByCategoryAndPeriod(ctx context.Context, categoryID, countryID, currencyID, month, year int) (int64, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	var result []struct {
		Total *int64 `json:"total"`
	}
	err := r.client.Transaction.Query().
		Where(
			transaction.HasCategoryWith(category.IDEQ(categoryID)),
			transaction.HasCountryWith(country.IDEQ(countryID)),
			transaction.HasAccountWith(account.HasCurrencyWith(currency.IDEQ(currencyID))),
			transaction.TransactionDateGTE(start),
			transaction.TransactionDateLT(end),
		).
		Aggregate(ent.As(ent.Sum(transaction.FieldAmount), "total")).
		Scan(ctx, &result)
	if err != nil {
		return 0, err
	}
	if len(result) == 0 || result[0].Total == nil {
		return 0, nil
	}
	return *result[0].Total, nil
}
