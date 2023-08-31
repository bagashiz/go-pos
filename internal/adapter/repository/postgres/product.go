package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/core/domain"
)

// CreateProduct creates a new product record in the database
func (db *DB) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	query := psql.Insert("products").
		Columns("category_id", "name", "image", "price", "stock").
		Values(product.CategoryID, product.Name, product.Image, product.Price, product.Stock).
		Suffix("RETURNING *").
		RunWith(db)

	err := query.QueryRowContext(ctx).Scan(
		&product.ID,
		&product.CategoryID,
		&product.SKU,
		&product.Name,
		&product.Stock,
		&product.Price,
		&product.Image,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID retrieves a product record from the database by id
func (db *DB) GetProductByID(ctx context.Context, id uint64) (*domain.Product, error) {
	query := psql.Select("*").
		From("products").
		Where(sq.Eq{"id": id}).
		Limit(1).
		RunWith(db)

	var product domain.Product

	err := query.QueryRowContext(ctx).Scan(
		&product.ID,
		&product.CategoryID,
		&product.SKU,
		&product.Name,
		&product.Stock,
		&product.Price,
		&product.Image,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

// ListProducts retrieves a list of products from the database
func (db *DB) ListProducts(ctx context.Context, search string, categoryId, skip, limit uint64) ([]*domain.Product, error) {
	query := psql.Select("*").
		From("products").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit).
		RunWith(db)

	if categoryId != 0 {
		query = query.Where(sq.Eq{"category_id": categoryId})
	}

	if search != "" {
		query = query.Where(sq.ILike{"name": "%" + search + "%"})
	}

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var products []*domain.Product

	for rows.Next() {
		var product domain.Product

		err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.SKU,
			&product.Name,
			&product.Stock,
			&product.Price,
			&product.Image,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

// UpdateProduct updates a product record in the database
func (db *DB) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	categoryId := nullUint64(product.CategoryID)
	name := nullString(product.Name)
	image := nullString(product.Image)
	price := nullFloat64(product.Price)
	stock := nullInt64(product.Stock)

	query := psql.Update("products").
		Set("name", sq.Expr("COALESCE(?, name)", name)).
		Set("category_id", sq.Expr("COALESCE(?, category_id)", categoryId)).
		Set("image", sq.Expr("COALESCE(?, image)", image)).
		Set("price", sq.Expr("COALESCE(?, price)", price)).
		Set("stock", sq.Expr("COALESCE(?, stock)", stock)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": product.ID}).
		Suffix("RETURNING *").
		RunWith(db)

	err := query.QueryRowContext(ctx).Scan(
		&product.ID,
		&product.CategoryID,
		&product.SKU,
		&product.Name,
		&product.Stock,
		&product.Price,
		&product.Image,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct deletes a product record from the database by id
func (db *DB) DeleteProduct(ctx context.Context, id uint64) error {
	query := psql.Delete("products").
		Where(sq.Eq{"id": id}).
		RunWith(db)

	_, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
