package controllers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"budgot/internal/ent"
	"budgot/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("fixed-dummy-value-for-timing-safety"), 12)

type UserFinder interface {
	FindByUsername(ctx context.Context, username string) (*ent.User, error)
}

type SessionCreator interface {
	Create(ctx context.Context, params repository.CreateSessionParams) (*ent.Session, error)
}

type AttemptRecorder interface {
	Record(ctx context.Context, username, ipAddress string, success bool) error
}

func LoginHandler(users UserFinder, sessions SessionCreator, attempts AttemptRecorder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		u, err := users.FindByUsername(r.Context(), username)
		if err != nil {
			bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
			_ = attempts.Record(r.Context(), username, r.RemoteAddr, false)
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
			_ = attempts.Record(r.Context(), username, r.RemoteAddr, false)
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		token := make([]byte, 32)
		if _, err := rand.Read(token); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		sum := sha256.Sum256(token)
		sessionID := hex.EncodeToString(sum[:])

		now := time.Now()

		sessionParams := repository.CreateSessionParams{
			ID:            sessionID,
			OwnerID:       u.ID,
			ExpiresAt:     now.Add(7 * 24 * time.Hour),
			LastSeen:      now,
			IPAddress:     r.RemoteAddr,
			UserAgentHash: hashUserAgent(r.UserAgent()),
		}

		_, err = sessions.Create(r.Context(), sessionParams)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    hex.EncodeToString(token),
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			Expires:  now.Add(7 * 24 * time.Hour),
		})

		_ = attempts.Record(r.Context(), username, r.RemoteAddr, true)
		w.Write([]byte("login successful"))
	}
}

func hashUserAgent(ua string) string {
	sum := sha256.Sum256([]byte(ua))
	return hex.EncodeToString(sum[:])
}
