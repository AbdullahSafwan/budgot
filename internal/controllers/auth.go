package controllers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"budgot/internal/ent"
	"budgot/internal/repository"

	chimw "github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/bcrypt"
)

type UserFinder interface {
	FindByUsername(ctx context.Context, username string) (*ent.User, error)
}

type SessionCreator interface {
	Create(ctx context.Context, params repository.CreateSessionParams) (*ent.Session, error)
}

type SessionSweeper interface {
	DeleteExpired(ctx context.Context, now time.Time) (int, error)
}

type AttemptRecorder interface {
	Record(ctx context.Context, username, ipAddress string, success bool) error
}

// AuthConfig holds auth-handler knobs sourced from configs.Config.
type AuthConfig struct {
	UseHTTPS   bool
	SessionTTL time.Duration
	BcryptCost int
}

func LoginHandler(users UserFinder, sessions SessionCreator, sweeper SessionSweeper, attempts AttemptRecorder, cfg AuthConfig) http.HandlerFunc {
	dummyHash, _ := bcrypt.GenerateFromPassword([]byte("fixed-dummy-value-for-timing-safety"), cfg.BcryptCost)

	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")
		clientIP := chimw.GetClientIP(r.Context())

		u, err := users.FindByUsername(r.Context(), username)
		if err != nil {
			bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
			recordAttempt(r.Context(), attempts, username, clientIP, false)
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
			recordAttempt(r.Context(), attempts, username, clientIP, false)
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
			ExpiresAt:     now.Add(cfg.SessionTTL),
			LastSeen:      now,
			IPAddress:     clientIP,
			UserAgentHash: hashUserAgent(r.UserAgent()),
		}

		_, err = sessions.Create(r.Context(), sessionParams)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if _, err := sweeper.DeleteExpired(r.Context(), now); err != nil {
			slog.Error("failed to sweep expired sessions", "error", err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    hex.EncodeToString(token),
			Path:     "/",
			HttpOnly: true,
			Secure:   cfg.UseHTTPS,
			SameSite: http.SameSiteLaxMode,
			Expires:  now.Add(cfg.SessionTTL),
		})

		recordAttempt(r.Context(), attempts, username, clientIP, true)
		w.Write([]byte("login successful"))
	}
}

func recordAttempt(ctx context.Context, attempts AttemptRecorder, username, ip string, success bool) {
	if err := attempts.Record(ctx, username, ip, success); err != nil {
		slog.Error("failed to record login attempt", "error", err, "username", username)
	}
}

func hashUserAgent(ua string) string {
	sum := sha256.Sum256([]byte(ua))
	return hex.EncodeToString(sum[:])
}

type SessionDeleter interface {
	Delete(ctx context.Context, id string) error
}

func LogoutHandler(sessions SessionDeleter, useHTTPS bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenBytes, err := hex.DecodeString(cookie.Value)
		if err == nil {
			sum := sha256.Sum256(tokenBytes)
			sessionID := hex.EncodeToString(sum[:])
			if err := sessions.Delete(r.Context(), sessionID); err != nil {
				slog.Error("failed to delete session on logout", "error", err)
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   useHTTPS,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
