package router

import (
	"net/http"
	"os"
	"text/template"
	"time"

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

func NewRouter(db *ent.Client) *chi.Mux {
	router := chi.NewRouter()
	router.Use(chimw.ClientIPFromXFFTrustedProxies(1))
	router.Use(chimw.Recoverer)
	router.Use(middleware.SecurityHeaders)

	csrfKey := []byte(os.Getenv("CSRF_AUTH_KEY"))

	//TODO update to true after deployment and tls
	router.Use(csrf.Protect(csrfKey, csrf.Secure(false)))

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

	var loginTmpl = template.Must(template.ParseFS(web.Templates, "templates/layout.html", "templates/login.html"))

	router.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		loginTmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
		})
	})

	router.Post("/logout", controllers.LogoutHandler(sessions))
	return router
}
