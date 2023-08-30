package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/core/domain"
)

/**
 * CategoryRepository implements port.CategoryRepository interface
 * and crovides an access to the postgres database
 */
type CategoryRepository struct {
	db *DB
}

// NewCategoryRepository creates a new category repository instance
func NewCategoryRepository(db *DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// CreateCategory creates a new category record in the database
func (cr *CategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	query := psql.Insert("categories").
		Columns("name").
		Values(category.Name).
		Suffix("RETURNING *").
		RunWith(cr.db)

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
func (cr *CategoryRepository) GetCategoryByID(ctx context.Context, id uint64) (*domain.Category, error) {
	query := psql.Select("*").
		From("categories").
		Where(sq.Eq{"id": id}).
		Limit(1).
		RunWith(cr.db)

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
func (cr *CategoryRepository) ListCategories(ctx context.Context, skip, limit uint64) ([]*domain.Category, error) {
	query := psql.Select("*").
		From("categories").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit).
		RunWith(cr.db)

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
func (cr *CategoryRepository) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	query := psql.Update("categories").
		Set("name", category.Name).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": category.ID}).
		Suffix("RETURNING *").
		RunWith(cr.db)

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
func (cr *CategoryRepository) DeleteCategory(ctx context.Context, id uint64) error {
	query := psql.Delete("categories").
		Where(sq.Eq{"id": id}).
		RunWith(cr.db)

	_, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
