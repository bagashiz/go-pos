package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

//go:generate mockgen -source=product.go -destination=mock/product.go -package=mock

// ProductRepository is an interface for interacting with product-related data
type ProductRepository interface {
	// CreateProduct inserts a new product into the database
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	// GetProductByID selects a product by id
	GetProductByID(ctx context.Context, id uint64) (*domain.Product, error)
	// ListProducts selects a list of products with pagination
	ListProducts(ctx context.Context, search string, categoryId, skip, limit uint64) ([]domain.Product, error)
	// UpdateProduct updates a product
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	// DeleteProduct deletes a product
	DeleteProduct(ctx context.Context, id uint64) error
}

// ProductService is an interface for interacting with product-related business logic
type ProductService interface {
	// CreateProduct creates a new product
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	// GetProduct returns a product by id
	GetProduct(ctx context.Context, id uint64) (*domain.Product, error)
	// ListProducts returns a list of products with pagination
	ListProducts(ctx context.Context, search string, categoryId, skip, limit uint64) ([]domain.Product, error)
	// UpdateProduct updates a product
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	// DeleteProduct deletes a product
	DeleteProduct(ctx context.Context, id uint64) error
}
