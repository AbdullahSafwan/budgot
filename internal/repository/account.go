package repository

import (
	"context"

	"budgot/internal/ent"
	"budgot/internal/ent/account"
	"budgot/internal/ent/user"
)

type AccountRepository struct {
	client *ent.Client
}

func NewAccountRepository(client *ent.Client) *AccountRepository {
	return &AccountRepository{client: client}
}

type CreateAccountParams struct {
	OwnerID     int
	CountryID   int
	CurrencyID  int
	Name        string
	AccountType account.AccountType
}

func (r *AccountRepository) Create(ctx context.Context, params CreateAccountParams) (*ent.Account, error) {
	return r.client.Account.Create().
		SetOwnerID(params.OwnerID).
		SetCountryID(params.CountryID).
		SetCurrencyID(params.CurrencyID).
		SetName(params.Name).
		SetAccountType(params.AccountType).
		Save(ctx)
}

func (r *AccountRepository) FindByID(ctx context.Context, id int) (*ent.Account, error) {
	return r.client.Account.Query().Where(account.IDEQ(id)).Only(ctx)
}

func (r *AccountRepository) ListByOwner(ctx context.Context, ownerID int) ([]*ent.Account, error) {
	return r.client.Account.Query().
		Where(account.HasOwnerWith(user.IDEQ(ownerID))).
		All(ctx)
}
