package service

import (
	"context"
	"errors"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
)

/**
 * CategoryService implements port.CategoryService interface
 * and provides an access to the category repository
 */
type CategoryService struct {
	repo port.CategoryRepository
}

// NewCategoryService creates a new category service instance
func NewCategoryService(repo port.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

// CreateCategory creates a new category
func (cs *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	return cs.repo.CreateCategory(ctx, category)
}

// GetCategory retrieves a category by id
func (cs *CategoryService) GetCategory(ctx context.Context, id uint64) (*domain.Category, error) {
	return cs.repo.GetCategoryByID(ctx, id)
}

// ListCategories retrieves a list of categories
func (cs *CategoryService) ListCategories(ctx context.Context, skip, limit uint64) ([]*domain.Category, error) {
	return cs.repo.ListCategories(ctx, skip, limit)
}

// UpdateCategory updates a category
func (cs *CategoryService) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	existingCategory, err := cs.repo.GetCategoryByID(ctx, category.ID)
	if err != nil {
		return nil, err
	}

	sameData := existingCategory.Name == category.Name
	if sameData {
		return nil, errors.New("no data to update")
	}

	return cs.repo.UpdateCategory(ctx, category)
}

// DeleteCategory deletes a category
func (cs *CategoryService) DeleteCategory(ctx context.Context, id uint64) error {
	_, err := cs.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return err
	}

	return cs.repo.DeleteCategory(ctx, id)
}
