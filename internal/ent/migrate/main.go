//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"budgot/internal/ent/migrate"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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
		schema.WithDialect(dialect.MySQL),
		schema.WithFormatter(atlas.DefaultFormatter),
	}

	if len(os.Args) != 2 {
		log.Fatalln("migration name required: go run -mod=mod internal/ent/migrate/main.go <name>")
	}

	dsn := fmt.Sprintf("mysql://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	if err := migrate.NamedDiff(ctx, dsn, os.Args[1], opts...); err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}
