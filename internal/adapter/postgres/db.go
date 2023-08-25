package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB is a wrapper for PostgreSQL database connection
type DB struct {
	*sql.DB
}

// NewDB creates a new PostgreSQL database instance
func NewDB() (*DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &DB{
		db,
	}, nil
}
