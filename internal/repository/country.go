package repository

import (
	"context"

	"budgot/internal/configs"
	"budgot/internal/ent"
	"budgot/internal/ent/country"
)

type CountryRepository struct {
	client *ent.Client
}

func NewCountryRepository(client *ent.Client) *CountryRepository {
	return &CountryRepository{client: client}
}

// WithTx binds the repository to an existing transaction.
func (r *CountryRepository) WithTx(tx *ent.Tx) *CountryRepository {
	return &CountryRepository{client: tx.Client()}
}

func (r *CountryRepository) Create(ctx context.Context, code, name string) (*ent.Country, error) {
	c, err := r.client.Country.Create().
		SetCode(code).
		SetName(name).
		Save(ctx)
	return c, configs.Translate(err)
}

func (r *CountryRepository) FindByCode(ctx context.Context, code string) (*ent.Country, error) {
	c, err := r.client.Country.Query().
		Where(country.CodeEQ(code)).
		WithDefaultCurrency().
		Only(ctx)
	return c, configs.Translate(err)
}

func (r *CountryRepository) List(ctx context.Context) ([]*ent.Country, error) {
	return r.client.Country.Query().WithDefaultCurrency().All(ctx)
}

// SetDefaultCurrency sets the currency to prefill in forms for new accounts
// and budgets in this country; it's a suggestion only, not enforced.
func (r *CountryRepository) SetDefaultCurrency(ctx context.Context, countryID, currencyID int) error {
	_, err := r.client.Country.UpdateOneID(countryID).SetDefaultCurrencyID(currencyID).Save(ctx)
	return configs.Translate(err)
}
