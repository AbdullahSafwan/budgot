package repository

import (
	"context"
	"time"

	"budgot/internal/ent"
)

type SessionRepository struct {
	client *ent.Client
}

func NewSessionRepository(client *ent.Client) *SessionRepository {
	return &SessionRepository{client: client}
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
