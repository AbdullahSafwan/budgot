package configs

import (
	"budgot/internal/ent"

	_ "modernc.org/sqlite"
)

func NewDB() (*ent.Client, error) {
	return ent.Open("sqlite3", "file:budgot.db?_fk=1&cache=shared")
}
