package controllers

import (
	"net/http"

	"budgot/internal/ent"
	"budgot/internal/ent/user"

	"golang.org/x/crypto/bcrypt"
)

var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("fixed-dummy-value-for-timing-safety"), bcrypt.DefaultCost)

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

		w.Write([]byte("login successful"))
	}
}
