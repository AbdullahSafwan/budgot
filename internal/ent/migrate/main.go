//go:build ignore

package main

import (
	"context"
	"log"
	"os"

	"budgot/internal/ent/migrate"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	dir, err := atlas.NewLocalDir("internal/ent/migrate/migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}

	opts := []schema.MigrateOption{
		schema.WithDir(dir),
		schema.WithMigrationMode(schema.ModeReplay),
		schema.WithDialect(dialect.SQLite),
		schema.WithFormatter(atlas.DefaultFormatter),
	}

	if len(os.Args) != 2 {
		log.Fatalln("migration name required: go run -mod=mod internal/ent/migrate/main.go <name>")
	}

	dsn := "sqlite://file::memory:?cache=shared&_fk=1"

	if err := migrate.NamedDiff(ctx, dsn, os.Args[1], opts...); err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}
