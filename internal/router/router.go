package router

import (
	"net/http"
	"time"

	"budgot/internal/controllers"
	"budgot/internal/ent"
	"budgot/internal/middleware"
	"budgot/internal/repository"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/httprate"
)

func NewRouter(db *ent.Client) *chi.Mux {
	router := chi.NewRouter()
	router.Use(chimw.ClientIPFromXFFTrustedProxies(1))

	loginLimiter := httprate.LimitBy(5, time.Minute, func(r *http.Request) (string, error) {
		return httprate.CanonicalizeIP(chimw.GetClientIP(r.Context())), nil
	})

	users := repository.NewUserRepository(db)
	sessions := repository.NewSessionRepository(db)
	attempts := repository.NewLoginAttemptRepository(db)

	router.With(middleware.RequireAuth(sessions)).Get("/", func(w http.ResponseWriter, r *http.Request) {
		u := middleware.UserFromContext(r.Context())
		w.Write([]byte("Welcome, " + u.Username))
	})

	router.With(loginLimiter).Post("/login", controllers.LoginHandler(users, sessions, attempts))
	router.With(loginLimiter).Get("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("login page placeholder"))
	})

	router.Post("/logout", controllers.LogoutHandler(sessions))
	return router
}
