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

// WithTx binds the repository to an existing transaction.
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
		Where(
			account.HasOwnerWith(user.IDEQ(ownerID)),
			account.IsActiveEQ(true),
		).
		All(ctx)
}

// Delete soft-deletes an account; rows are never hard-deleted.
func (r *AccountRepository) Delete(ctx context.Context, id int) error {
	_, err := r.client.Account.UpdateOneID(id).SetIsActive(false).Save(ctx)
	return err
}
