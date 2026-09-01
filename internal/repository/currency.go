package repository

import (
	"context"

	"budgot/internal/configs"
	"budgot/internal/ent"
	"budgot/internal/ent/currency"
)

type CurrencyRepository struct {
	client *ent.Client
}

func NewCurrencyRepository(client *ent.Client) *CurrencyRepository {
	return &CurrencyRepository{client: client}
}

// WithTx binds the repository to an existing transaction.
func (r *CurrencyRepository) WithTx(tx *ent.Tx) *CurrencyRepository {
	return &CurrencyRepository{client: tx.Client()}
}

func (r *CurrencyRepository) Create(ctx context.Context, code, name, symbol string, decimalPlaces int) (*ent.Currency, error) {
	c, err := r.client.Currency.Create().
		SetCode(code).
		SetName(name).
		SetSymbol(symbol).
		SetDecimalPlaces(decimalPlaces).
		Save(ctx)
	return c, configs.Translate(err)
}

func (r *CurrencyRepository) FindByCode(ctx context.Context, code string) (*ent.Currency, error) {
	c, err := r.client.Currency.Query().Where(currency.CodeEQ(code)).Only(ctx)
	return c, configs.Translate(err)
}

func (r *CurrencyRepository) List(ctx context.Context) ([]*ent.Currency, error) {
	return r.client.Currency.Query().All(ctx)
}
