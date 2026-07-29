package configs

import (
	"fmt"
	"os"

	"budgot/internal/ent"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() (*ent.Client, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	return ent.Open("mysql", dsn)

}
