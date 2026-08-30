package repository

import (
	"context"
	"time"

	"budgot/internal/ent"
	"budgot/internal/ent/session"
)

type SessionRepository struct {
	client *ent.Client
}

func NewSessionRepository(client *ent.Client) *SessionRepository {
	return &SessionRepository{client: client}
}

// WithTx binds the repository to an existing transaction.
func (r *SessionRepository) WithTx(tx *ent.Tx) *SessionRepository {
	return &SessionRepository{client: tx.Client()}
}

type CreateSessionParams struct {
	ID            string
	OwnerID       int
	ExpiresAt     time.Time
	LastSeen      time.Time
	IPAddress     string
	UserAgentHash string
}

func (r *SessionRepository) Create(ctx context.Context, params CreateSessionParams) (*ent.Session, error) {
	return r.client.Session.Create().
		SetID(params.ID).
		SetOwnerID(params.OwnerID).
		SetExpiresAt(params.ExpiresAt).
		SetLastSeen(params.LastSeen).
		SetIPAddress(params.IPAddress).
		SetUserAgentHash(params.UserAgentHash).
		Save(ctx)
}

func (r *SessionRepository) FindByID(ctx context.Context, id string) (*ent.Session, error) {
	return r.client.Session.Query().
		Where(session.IDEQ(id)).
		WithOwner().Only(ctx)
}

func (r *SessionRepository) UpdateLastSeen(ctx context.Context, id string, t time.Time) error {
	_, err := r.client.Session.UpdateOneID(id).SetLastSeen(t).Save(ctx)
	return err
}

func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	return r.client.Session.DeleteOneID(id).Exec(ctx)
}

// DeleteExpired hard-deletes sessions past their expiry (the one exception to soft-delete).
func (r *SessionRepository) DeleteExpired(ctx context.Context, now time.Time) (int, error) {
	return r.client.Session.Delete().
		Where(session.ExpiresAtLT(now)).
		Exec(ctx)
}
