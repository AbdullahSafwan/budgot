package repository

import (
	"context"

	"budgot/internal/ent"
	"budgot/internal/ent/country"
)

type CountryRepository struct {
	client *ent.Client
}

func NewCountryRepository(client *ent.Client) *CountryRepository {
	return &CountryRepository{client: client}
}

// WithTx returns a copy of the repository bound to the given transaction, so its
func (r *CountryRepository) WithTx(tx *ent.Tx) *CountryRepository {
	return &CountryRepository{client: tx.Client()}
}

func (r *CountryRepository) Create(ctx context.Context, code, name string) (*ent.Country, error) {
	return r.client.Country.Create().
		SetCode(code).
		SetName(name).
		Save(ctx)
}

func (r *CountryRepository) FindByCode(ctx context.Context, code string) (*ent.Country, error) {
	return r.client.Country.Query().Where(country.CodeEQ(code)).Only(ctx)
}

func (r *CountryRepository) List(ctx context.Context) ([]*ent.Country, error) {
	return r.client.Country.Query().All(ctx)
}
