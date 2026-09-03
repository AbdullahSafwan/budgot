package middleware

import (
	"budgot/internal/ent"
	"context"
	"net/http"
)

const countryContextKey contextKey = "country"

type CountryStore interface {
	FindByCode(ctx context.Context, code string) (*ent.Country, error)
}

// RequireCountry resolves ?country=<code> into context; run after RequireAuth.
func RequireCountry(countries CountryStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("country")
			if code == "" {
				http.Error(w, "country is required", http.StatusBadRequest)
				return
			}

			c, err := countries.FindByCode(r.Context(), code)
			if err != nil {
				http.Error(w, "unknown country", http.StatusNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), countryContextKey, c)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CountryFromContext(ctx context.Context) *ent.Country {
	c, _ := ctx.Value(countryContextKey).(*ent.Country)
	return c
}
