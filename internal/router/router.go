package router

import (
	"net/http"
	"time"

	"budgot/internal/controllers"
	"budgot/internal/ent"
	"budgot/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/httprate"
)

func NewRouter(db *ent.Client) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.ClientIPFromXFFTrustedProxies(1))

	loginLimiter := httprate.LimitBy(5, time.Minute, func(r *http.Request) (string, error) {
		return httprate.CanonicalizeIP(middleware.GetClientIP(r.Context())), nil
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Budgot API!"))
	})

	users := repository.NewUserRepository(db)
	sessions := repository.NewSessionRepository(db)
	attempts := repository.NewLoginAttemptRepository(db)

	router.With(loginLimiter).Post("/login", controllers.LoginHandler(users, sessions, attempts))
	return router
}
