package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/core/domain"
)

// CreateCategory creates a new category record in the database
func (db *DB) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	query := psql.Insert("categories").
		Columns("name").
		Values(category.Name).
		Suffix("RETURNING *").
		RunWith(db)

	err := query.QueryRowContext(ctx).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID retrieves a category record from the database by id
func (db *DB) GetCategoryByID(ctx context.Context, id uint64) (*domain.Category, error) {
	query := psql.Select("*").
		From("categories").
		Where(sq.Eq{"id": id}).
		Limit(1).
		RunWith(db)

	var category domain.Category

	err := query.QueryRowContext(ctx).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &category, nil
}

// ListCategories retrieves a list of categories from the database
func (db *DB) ListCategories(ctx context.Context, skip, limit uint64) ([]*domain.Category, error) {
	query := psql.Select("*").
		From("categories").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit).
		RunWith(db)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var categories []*domain.Category

	for rows.Next() {
		var category domain.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

// UpdateCategory updates a category record in the database
func (db *DB) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	query := psql.Update("categories").
		Set("name", category.Name).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": category.ID}).
		Suffix("RETURNING *").
		RunWith(db)

	err := query.QueryRowContext(ctx).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category record from the database by id
func (db *DB) DeleteCategory(ctx context.Context, id uint64) error {
	query := psql.Delete("categories").
		Where(sq.Eq{"id": id}).
		RunWith(db)

	_, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
