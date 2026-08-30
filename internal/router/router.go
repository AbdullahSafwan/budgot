package router

import (
	"html/template"
	"log/slog"
	"net/http"

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

// app holds the dependencies handlers are built from.
type app struct {
	cfg       configs.Config
	users     *repository.UserRepository
	sessions  *repository.SessionRepository
	attempts  *repository.LoginAttemptRepository
	loginTmpl *template.Template
}

func NewRouter(db *ent.Client, cfg configs.Config) *chi.Mux {
	a := &app{
		cfg:       cfg,
		users:     repository.NewUserRepository(db),
		sessions:  repository.NewSessionRepository(db),
		attempts:  repository.NewLoginAttemptRepository(db),
		loginTmpl: template.Must(template.ParseFS(web.Templates, "templates/layout.html", "templates/login.html")),
	}

	r := chi.NewRouter()
	r.NotFound(notFoundHandler)
	r.MethodNotAllowed(methodNotAllowedHandler)

	// Trust XFF only when a proxy is configured; otherwise use the raw TCP peer.
	if cfg.TrustedProxyHops > 0 {
		r.Use(chimw.ClientIPFromXFFTrustedProxies(cfg.TrustedProxyHops))
	} else {
		r.Use(chimw.ClientIPFromRemoteAddr)
	}
	r.Use(middleware.RequestLogger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.SecurityHeaders(cfg.IsProduction()))

	csrfKey := []byte(cfg.CSRFAuthKey)

	if !cfg.IsProduction() {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, csrf.PlaintextHTTPRequest(r))
			})
		})
	}

	r.Use(csrf.Protect(csrfKey,
		csrf.Secure(cfg.IsProduction()),
		csrf.ErrorHandler(http.HandlerFunc(csrfErrorHandler)),
	))

	loginLimiter := httprate.LimitBy(cfg.LoginRateLimit, cfg.LoginRateLimitWindow, func(r *http.Request) (string, error) {
		return httprate.CanonicalizeIP(chimw.GetClientIP(r.Context())), nil
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(a.sessions, cfg.SessionIdleTimeout))
		r.Get("/", a.dashboard)
	})

	r.With(loginLimiter).Post("/login", a.loginHandler())
	r.Get("/login", a.loginPage)
	r.Post("/logout", a.logoutHandler())

	return r
}

func (a *app) dashboard(w http.ResponseWriter, r *http.Request) {
	u := middleware.UserFromContext(r.Context())
	w.Write([]byte("Welcome, " + u.Username))
}

func (a *app) loginHandler() http.HandlerFunc {
	authCfg := controllers.AuthConfig{
		UseHTTPS:   a.cfg.IsProduction(),
		SessionTTL: a.cfg.SessionTTL,
		BcryptCost: a.cfg.BcryptCost,
	}
	return controllers.LoginHandler(a.users, a.sessions, a.sessions, a.attempts, authCfg)
}

func (a *app) logoutHandler() http.HandlerFunc {
	return controllers.LogoutHandler(a.sessions, a.cfg.IsProduction())
}

func (a *app) loginPage(w http.ResponseWriter, r *http.Request) {
	if err := a.loginTmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}); err != nil {
		slog.Error("failed to render login page", "error", err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not found", http.StatusNotFound)
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func csrfErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "form expired or invalid — please refresh and try again", http.StatusForbidden)
}
