package repository

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

/**
 * ProductRepository implements port.ProductRepository interface
 * and provides an access to the postgres database
 */
type ProductRepository struct {
	db *DB
}

// NewProductRepository creates a new product repository instance
func NewProductRepository(db *DB) *ProductRepository {
	return &ProductRepository{
		db,
	}
}

// CreateProduct creates a new product record in the database
func (pr *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	query := psql.Insert("products").
		Columns("category_id", "name", "image", "price", "stock").
		Values(product.CategoryID, product.Name, product.Image, product.Price, product.Stock).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.db.QueryRow(ctx, sql, args...).Scan(
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
func (pr *ProductRepository) GetProductByID(ctx context.Context, id uint64) (*domain.Product, error) {
	var product domain.Product

	query := psql.Select("*").
		From("products").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.db.QueryRow(ctx, sql, args...).Scan(
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
		if err == pgx.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

// ListProducts retrieves a list of products from the database
func (pr *ProductRepository) ListProducts(ctx context.Context, search string, categoryId, skip, limit uint64) ([]*domain.Product, error) {
	var product domain.Product
	var products []*domain.Product

	query := psql.Select("*").
		From("products").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	if categoryId != 0 {
		query = query.Where(sq.Eq{"category_id": categoryId})
	}

	if search != "" {
		query = query.Where(sq.ILike{"name": "%" + search + "%"})
	}

	rows, err := pr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
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
func (pr *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
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
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.db.QueryRow(ctx, sql, args...).Scan(
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
func (pr *ProductRepository) DeleteProduct(ctx context.Context, id uint64) error {
	query := psql.Delete("products").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = pr.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
