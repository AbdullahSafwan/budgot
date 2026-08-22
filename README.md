# budgot

A personal budget tracking app — single Go binary, server-rendered HTML (`html/template`), no SPA. Built for personal, multi-country, multi-currency use.

Full design details, schema, and the reasoning behind every architectural decision live in [budget-app-specification.md](budget-app-specification.md) — this file is just the quick-start.

## Status

- **Authentication** — done. Login, logout, session management, CSRF, rate limiting, audit logging, security headers.
- **Budgeting** (accounts, categories, transactions, budgets) — designed, not yet implemented.

## Tech stack

Go, chi router, SQLite (`modernc.org/sqlite`, no CGO), Ent ORM with Atlas-managed migrations, `gorilla/csrf`, `go-chi/httprate`, `html/template` served via `go:embed`.

## Setup

Requires a `.env` file at the repo root with:

```
CSRF_AUTH_KEY=   # 32 bytes, e.g. `openssl rand -base64 32`
PORT=            # e.g. 8080
APP_ENV=         # "development" or "production"
```

All three are required — the app fails fast on startup if any are missing, rather than silently defaulting.

## Commands

```
go build ./...              # build
go run ./cmd/server          # run
go vet ./...                 # vet
```

### Migrations (Atlas + Ent)

Migration files live in `internal/ent/migrate/migrations/`.

```
go generate ./...                                                  # regenerate Ent code after a schema change
go run -mod=mod internal/ent/migrate/main.go <migration_name>      # generate a new migration
set -a; source .env; set +a; atlas migrate apply --env local        # apply pending migrations
```

Ent schema source lives in `internal/ent/schema/` — that's the only hand-edited part of `internal/ent/`; everything else is generated.
