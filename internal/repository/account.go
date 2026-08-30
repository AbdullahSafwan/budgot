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

// WithTx returns a copy of// WithTx returns a copy of the repository bound to the given transaction, so its
func (r *AccountRepository) WithTx(tx *ent.Tx) *AccountRepository {
	return &AccountRepository{client: tx.Client()}
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
