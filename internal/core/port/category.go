package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

// CategoryRepository is an interface for interacting with category-related data
type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	GetCategoryByID(ctx context.Context, id uint64) (*domain.Category, error)
	ListCategories(ctx context.Context, skip, limit uint64) ([]*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
}

// CategoryService is an interface for interacting with category-related business logic
type CategoryService interface {
	CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	GetCategory(ctx context.Context, id uint64) (*domain.Category, error)
	ListCategories(ctx context.Context, skip, limit uint64) ([]*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
}
