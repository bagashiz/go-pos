package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

//go:generate mockgen -source=category.go -destination=mock/category.go -package=mock

// CategoryRepository is an interface for interacting with category-related data
type CategoryRepository interface {
	// CreateCategory inserts a new category into the database
	CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	// GetCategoryByID selects a category by id
	GetCategoryByID(ctx context.Context, id uint64) (*domain.Category, error)
	// ListCategories selects a list of categories with pagination
	ListCategories(ctx context.Context, skip, limit uint64) ([]domain.Category, error)
	// UpdateCategory updates a category
	UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	// DeleteCategory deletes a category
	DeleteCategory(ctx context.Context, id uint64) error
}

// CategoryService is an interface for interacting with category-related business logic
type CategoryService interface {
	// CreateCategory creates a new category
	CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	// GetCategory returns a category by id
	GetCategory(ctx context.Context, id uint64) (*domain.Category, error)
	// ListCategories returns a list of categories with pagination
	ListCategories(ctx context.Context, skip, limit uint64) ([]domain.Category, error)
	// UpdateCategory updates a category
	UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	// DeleteCategory deletes a category
	DeleteCategory(ctx context.Context, id uint64) error
}
