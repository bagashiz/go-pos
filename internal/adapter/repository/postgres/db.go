package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

/**
 * DB is a wrapper for PostgreSQL database connection
 * that uses pgxpool as database driver
 */
type DB struct {
	*pgxpool.Pool
}

// NewDB creates a new PostgreSQL database instance
func NewDB(ctx context.Context) (*DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &DB{
		db,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	db.Pool.Close()
}
