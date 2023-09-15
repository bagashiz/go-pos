package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

// ProductRepository is an interface for interacting with product-related data
type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProductByID(ctx context.Context, id uint64) (*domain.Product, error)
	ListProducts(ctx context.Context, search string, categoryId, skip, limit uint64) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id uint64) error
}

// ProductService is an interface for interacting with product-related business logic
type ProductService interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProduct(ctx context.Context, id uint64) (*domain.Product, error)
	ListProducts(ctx context.Context, search string, categoryId, skip, limit uint64) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id uint64) error
}
