package service

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/bagashiz/go-pos/internal/core/util"
)

/**
 * CategoryService implements port.CategoryService interface
 * and provides an access to the category repository
 * and cache service
 */
type CategoryService struct {
	repo  port.CategoryRepository
	cache port.CacheService
}

// NewCategoryService creates a new category service instance
func NewCategoryService(repo port.CategoryRepository, cache port.CacheService) *CategoryService {
	return &CategoryService{
		repo,
		cache,
	}
}

// CreateCategory creates a new category
func (cs *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	_, err := cs.repo.CreateCategory(ctx, category)
	if err != nil {
		if port.IsUniqueConstraintViolationError(err) {
			return nil, port.ErrConflictingData
		}

		return nil, err
	}

	cacheKey := util.GenerateCacheKey("category", category.ID)
	categorySerialized, err := util.Serialize(category)
	if err != nil {
		return nil, err
	}

	err = cs.cache.Set(ctx, cacheKey, categorySerialized, 0)
	if err != nil {
		return nil, err
	}

	err = cs.cache.DeleteByPrefix(ctx, "categories:*")
	if err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategory retrieves a category by id
func (cs *CategoryService) GetCategory(ctx context.Context, id uint64) (*domain.Category, error) {
	var category *domain.Category

	cacheKey := util.GenerateCacheKey("category", id)
	cachedCategory, err := cs.cache.Get(ctx, cacheKey)
	if err == nil {
		err := util.Deserialize(cachedCategory, &category)
		if err != nil {
			return nil, err
		}

		return category, nil
	}

	category, err = cs.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	categorySerialized, err := util.Serialize(category)
	if err != nil {
		return nil, err
	}

	err = cs.cache.Set(ctx, cacheKey, categorySerialized, 0)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// ListCategories retrieves a list of categories
func (cs *CategoryService) ListCategories(ctx context.Context, skip, limit uint64) ([]domain.Category, error) {
	var categories []domain.Category

	params := util.GenerateCacheKeyParams(skip, limit)
	cacheKey := util.GenerateCacheKey("categories", params)

	cachedCategories, err := cs.cache.Get(ctx, cacheKey)
	if err == nil {
		err := util.Deserialize(cachedCategories, &categories)
		if err != nil {
			return nil, err
		}

		return categories, nil
	}

	categories, err = cs.repo.ListCategories(ctx, skip, limit)
	if err != nil {
		return nil, err
	}

	categoriesSerialized, err := util.Serialize(categories)
	if err != nil {
		return nil, err
	}

	err = cs.cache.Set(ctx, cacheKey, categoriesSerialized, 0)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// UpdateCategory updates a category
func (cs *CategoryService) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	existingCategory, err := cs.repo.GetCategoryByID(ctx, category.ID)
	if err != nil {
		return nil, err
	}

	emptyData := category.Name == ""
	sameData := existingCategory.Name == category.Name
	if emptyData || sameData {
		return nil, port.ErrNoUpdatedData
	}

	_, err = cs.repo.UpdateCategory(ctx, category)
	if err != nil {
		if port.IsUniqueConstraintViolationError(err) {
			return nil, port.ErrConflictingData
		}

		return nil, err
	}

	cacheKey := util.GenerateCacheKey("category", category.ID)
	_ = cs.cache.Delete(ctx, cacheKey)

	categorySerialized, err := util.Serialize(category)
	if err != nil {
		return nil, err
	}

	err = cs.cache.Set(ctx, cacheKey, categorySerialized, 0)
	if err != nil {
		return nil, err
	}

	err = cs.cache.DeleteByPrefix(ctx, "categories:*")
	if err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category
func (cs *CategoryService) DeleteCategory(ctx context.Context, id uint64) error {
	_, err := cs.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return err
	}

	cacheKey := util.GenerateCacheKey("category", id)
	_ = cs.cache.Delete(ctx, cacheKey)

	err = cs.cache.DeleteByPrefix(ctx, "categories:*")
	if err != nil {
		return err
	}

	return cs.repo.DeleteCategory(ctx, id)
}
