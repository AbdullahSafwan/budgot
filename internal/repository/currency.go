package repository

import (
	"context"

	"budgot/internal/ent"
	"budgot/internal/ent/currency"
)

type CurrencyRepository struct {
	client *ent.Client
}

func NewCurrencyRepository(client *ent.Client) *CurrencyRepository {
	return &CurrencyRepository{client: client}
}

// WithTx returns a copy of the repository bound to the given transaction, so its
func (r *CurrencyRepository) WithTx(tx *ent.Tx) *CurrencyRepository {
	return &CurrencyRepository{client: tx.Client()}
}

func (r *CurrencyRepository) Create(ctx context.Context, code, name, symbol string) (*ent.Currency, error) {
	return r.client.Currency.Create().
		SetCode(code).
		SetName(name).
		SetSymbol(symbol).
		Save(ctx)
}

func (r *CurrencyRepository) FindByCode(ctx context.Context, code string) (*ent.Currency, error) {
	return r.client.Currency.Query().Where(currency.CodeEQ(code)).Only(ctx)
}

func (r *CurrencyRepository) List(ctx context.Context) ([]*ent.Currency, error) {
	return r.client.Currency.Query().All(ctx)
}
