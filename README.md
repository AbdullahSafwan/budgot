# budgot

A personal budget tracking app — single Go binary, server-rendered HTML (`html/template`), no SPA. Built for personal, multi-country, multi-currency use.

## Goals

- Server-rendered HTML, clean layered architecture, easy to extend

## Tech stack

Go, chi router, SQLite (`modernc.org/sqlite`, no CGO), Ent ORM with Atlas-managed migrations, `gorilla/csrf`, `go-chi/httprate`, `html/template` served via `go:embed`, bcrypt.

## How it works

Budgot is **country-first**: each country is effectively a separate profile, with its own accounts, transactions, and budgets. A combined cross-country view exists, but it's secondary — the primary UX is picking a country and seeing only that country's data.

Currencies are **siloed** — no FX conversion, no cross-currency netting. An account belongs to one country and one currency, chosen independently (a country's accounts can span currencies); a budget is scoped by category, country, *and* currency for the same reason. Categories are the one thing shared across countries, since "Groceries" or "Rent" mean the same thing everywhere.

A transaction's `amount` is **signed**, independent of its category: negative means money left the account, positive means it arrived. A transfer between two accounts is just two transactions — one leg per account, opposite signs, linked by a shared `transfer_group` — rather than full double-entry bookkeeping.

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
