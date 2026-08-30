package middleware

import (
	"budgot/internal/ent"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

type contextKey string

const userContextKey contextKey = "user"

type SessionStore interface {
	FindByID(ctx context.Context, id string) (*ent.Session, error)
	UpdateLastSeen(ctx context.Context, id string, t time.Time) error
	Delete(ctx context.Context, id string) error
}

func RequireAuth(sessions SessionStore, idleTimeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			tokenBytes, err := hex.DecodeString(cookie.Value)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			sum := sha256.Sum256(tokenBytes)
			sessionID := hex.EncodeToString(sum[:])

			sess, err := sessions.FindByID(r.Context(), sessionID)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			now := time.Now()
			expired := now.After(sess.ExpiresAt) || now.After(sess.LastSeen.Add(idleTimeout))
			hijacked := sess.UserAgentHash != hashUserAgent(r.UserAgent())
			if expired || hijacked {
				if err := sessions.Delete(r.Context(), sess.ID); err != nil {
					slog.Error("failed to delete session", "error", err)
				}
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			if err := sessions.UpdateLastSeen(r.Context(), sess.ID, now); err != nil {
				slog.Error("failed to update session last-seen", "error", err)
			}

			ctx := context.WithValue(r.Context(), userContextKey, sess.Edges.Owner)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func hashUserAgent(ua string) string {
	sum := sha256.Sum256([]byte(ua))
	return hex.EncodeToString(sum[:])
}

func UserFromContext(ctx context.Context) *ent.User {
	u, _ := ctx.Value(userContextKey).(*ent.User)
	return u
}
