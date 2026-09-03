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
	"budgot/internal/ent/predicate"
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

func (r *TransactionRepository) FindByID(ctx context.Context, ownerID, id int) (*ent.Transaction, error) {
	t, err := r.client.Transaction.Query().
		Where(transaction.IDEQ(id), transaction.HasOwnerWith(user.IDEQ(ownerID))).
		Only(ctx)
	return t, configs.Translate(err)
}

type ListTransactionsParams struct {
	OwnerID       int
	CountryID     int
	AccountID     *int
	CategoryID    *int
	From, To      *time.Time
	WithEdges     bool
	Limit, Offset int
}

func (r *TransactionRepository) List(ctx context.Context, p ListTransactionsParams) ([]*ent.Transaction, error) {
	preds := []predicate.Transaction{
		transaction.HasOwnerWith(user.IDEQ(p.OwnerID)),
		transaction.HasCountryWith(country.IDEQ(p.CountryID)),
	}
	if p.AccountID != nil {
		preds = append(preds, transaction.HasAccountWith(account.IDEQ(*p.AccountID)))
	}
	if p.CategoryID != nil {
		preds = append(preds, transaction.HasCategoryWith(category.IDEQ(*p.CategoryID)))
	}
	if p.From != nil {
		preds = append(preds, transaction.TransactionDateGTE(*p.From))
	}
	if p.To != nil {
		preds = append(preds, transaction.TransactionDateLT(*p.To))
	}

	q := r.client.Transaction.Query().
		Where(preds...).
		Order(ent.Desc(transaction.FieldTransactionDate)).
		Limit(listLimit(p.Limit)).
		Offset(p.Offset)
	if p.WithEdges {
		q = q.WithAccount().WithCategory()
	}
	return q.All(ctx)
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

type SumByCategoryAndPeriodParams struct {
	CategoryID int
	CountryID  int
	CurrencyID int
	Month      int
	Year       int
}

// SumByCategoryAndPeriod sums a category's transactions for one month, filtered by country and currency.
func (r *TransactionRepository) SumByCategoryAndPeriod(ctx context.Context, p SumByCategoryAndPeriodParams) (int64, error) {
	start := time.Date(p.Year, time.Month(p.Month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	var result []struct {
		Total *int64 `json:"total"`
	}
	err := r.client.Transaction.Query().
		Where(
			transaction.HasCategoryWith(category.IDEQ(p.CategoryID)),
			transaction.HasCountryWith(country.IDEQ(p.CountryID)),
			transaction.HasAccountWith(account.HasCurrencyWith(currency.IDEQ(p.CurrencyID))),
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
