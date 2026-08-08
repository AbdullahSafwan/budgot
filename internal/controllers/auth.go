package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"budgot/internal/ent"
	"budgot/internal/ent/user"

	"golang.org/x/crypto/bcrypt"
)

var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("fixed-dummy-value-for-timing-safety"), 12)

func LoginHandler(db *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		u, err := db.User.Query().Where(user.UsernameEQ(username)).Only(r.Context())
		if err != nil {
			bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
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
		_, err = db.Session.Create().
			SetID(sessionID).
			SetOwnerID(u.ID).
			SetExpiresAt(now.Add(7 * 24 * time.Hour)).
			SetLastSeen(now).
			SetIPAddress(r.RemoteAddr).
			SetUserAgentHash(hashUserAgent(r.UserAgent())).
			Save(r.Context())
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

		w.Write([]byte("login successful"))
	}
}

func hashUserAgent(ua string) string {
	sum := sha256.Sum256([]byte(ua))
	return hex.EncodeToString(sum[:])
}
