package configs

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"budgot/internal/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"modernc.org/sqlite"
)

// Not "sqlite3" — that name collides with mattn/go-sqlite3 and would panic on registration.
const driverName = "sqlite3_budgot"

type sqliteDriver struct {
	*sqlite.Driver
}

func (d sqliteDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return conn, err
	}
	c := conn.(interface {
		Exec(stmt string, args []driver.Value) (driver.Result, error)
	})
	for _, pragma := range []string{
		"PRAGMA foreign_keys = on;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA busy_timeout = 5000;",
	} {
		if _, err := c.Exec(pragma, nil); err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to apply %q: %w", pragma, err)
		}
	}
	return conn, nil
}

func init() {
	sql.Register(driverName, sqliteDriver{Driver: &sqlite.Driver{}})
}

func NewDB(cfg Config) (*ent.Client, error) {
	db, err := sql.Open(driverName, cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	// SQLite serializes writes anyway; one conn removes lock contention entirely.
	db.SetMaxOpenConns(1)

	drv := entsql.OpenDB(dialect.SQLite, db)
	return ent.NewClient(ent.Driver(drv)), nil
}
