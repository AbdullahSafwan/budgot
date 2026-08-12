package middleware

import (
	"budgot/internal/ent"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"
)

type contextKey string

const userContextKey contextKey = "user"
const idleTimeout = 30 * time.Minute

type SessionStore interface {
	FindByID(ctx context.Context, id string) (*ent.Session, error)
	UpdateLastSeen(ctx context.Context, id string, t time.Time) error
	Delete(ctx context.Context, id string) error
}

func RequireAuth(sessions SessionStore) func(http.Handler) http.Handler {
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
			if now.After(sess.ExpiresAt) || now.After(sess.LastSeen.Add(idleTimeout)) {
				_ = sessions.Delete(r.Context(), sess.ID)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			_ = sessions.UpdateLastSeen(r.Context(), sess.ID, now)

			ctx := context.WithValue(r.Context(), userContextKey, sess.Edges.Owner)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) *ent.User {
	u, _ := ctx.Value(userContextKey).(*ent.User)
	return u
}
