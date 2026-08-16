package router

import (
	"net/http"
	"text/template"
	"time"

	"budgot/internal/configs"
	"budgot/internal/controllers"
	"budgot/internal/ent"
	"budgot/internal/middleware"
	"budgot/internal/repository"
	"budgot/web"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"

	"github.com/go-chi/httprate"
)

func NewRouter(db *ent.Client, cfg configs.Config) *chi.Mux {
	router := chi.NewRouter()
	router.Use(chimw.ClientIPFromXFFTrustedProxies(1))
	router.Use(chimw.Recoverer)
	router.Use(middleware.SecurityHeaders(cfg.IsProduction()))

	csrfKey := []byte(cfg.CSRFAuthKey)

	if !cfg.IsProduction() {
		//TODO Add logger
		// router.Use(middleware.DevLogger)
		router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, csrf.PlaintextHTTPRequest(r))
			})
		})
	}

	router.Use(csrf.Protect(csrfKey, csrf.Secure(cfg.IsProduction())))

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

	router.With(loginLimiter).Post("/login", controllers.LoginHandler(users, sessions, attempts, cfg.IsProduction()))

	var loginTmpl = template.Must(template.ParseFS(web.Templates, "templates/layout.html", "templates/login.html"))

	router.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		loginTmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
		})
	})

	router.Post("/logout", controllers.LogoutHandler(sessions, cfg.IsProduction()))
	return router
}
