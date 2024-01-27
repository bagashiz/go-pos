package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port/mock"
	"github.com/bagashiz/go-pos/internal/core/service"
	"github.com/bagashiz/go-pos/internal/core/util"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type createCategoryTestedInput struct {
	category *domain.Category
}

type createCategoryExpectedOutput struct {
	category *domain.Category
	err      error
}

func TestCategoryService_CreateCategory(t *testing.T) {
	ctx := context.Background()
	categoryID := gofakeit.Uint64()
	categoryName := gofakeit.ProductCategory()
	categoryInput := &domain.Category{
		Name: categoryName,
	}
	categoryOutput := &domain.Category{
		ID:        categoryID,
		Name:      categoryName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cacheKey := util.GenerateCacheKey("category", categoryOutput.ID)
	categorySerialized, _ := util.Serialize(categoryOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			categoryRepo *mock.MockCategoryRepository,
			cache *mock.MockCacheRepository,
		)
		input    createCategoryTestedInput
		expected createCategoryExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					CreateCategory(gomock.Any(), gomock.Eq(categoryInput)).
					Times(1).
					Return(categoryOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(categorySerialized), gomock.Eq(ttl)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("categories:*")).
					Times(1).
					Return(nil)
			},
			input: createCategoryTestedInput{
				category: categoryInput,
			},
			expected: createCategoryExpectedOutput{
				category: categoryOutput,
				err:      nil,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					CreateCategory(gomock.Any(), gomock.Eq(categoryInput)).
					Times(1).
					Return(nil, domain.ErrConflictingData)
			},
			input: createCategoryTestedInput{
				category: categoryInput,
			},
			expected: createCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrConflictingData,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					CreateCategory(gomock.Any(), gomock.Eq(categoryInput)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: createCategoryTestedInput{
				category: categoryInput,
			},
			expected: createCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					CreateCategory(gomock.Any(), gomock.Eq(categoryInput)).
					Times(1).
					Return(categoryOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(categorySerialized), gomock.Eq(ttl)).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: createCategoryTestedInput{
				category: categoryInput,
			},
			expected: createCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					CreateCategory(gomock.Any(), gomock.Eq(categoryInput)).
					Times(1).
					Return(categoryOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(categorySerialized), gomock.Eq(ttl)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("categories:*")).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: createCategoryTestedInput{
				category: categoryInput,
			},
			expected: createCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			categoryRepo := mock.NewMockCategoryRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(categoryRepo, cache)

			categoryService := service.NewCategoryService(categoryRepo, cache)

			category, err := categoryService.CreateCategory(ctx, tc.input.category)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.category, category, "Category mismatch")
		})
	}
}

type getCategoryTestedInput struct {
	id uint64
}

type getCategoryExpectedOutput struct {
	category *domain.Category
	err      error
}

func TestCategoryService_GetCategory(t *testing.T) {
	ctx := context.Background()
	categoryID := gofakeit.Uint64()
	categoryName := gofakeit.ProductCategory()
	category := &domain.Category{
		ID:   categoryID,
		Name: categoryName,
	}

	cacheKey := util.GenerateCacheKey("category", category.ID)
	categorySerialized, _ := util.Serialize(category)

	testCases := []struct {
		desc  string
		mocks func(
			categoryRepo *mock.MockCategoryRepository,
			cache *mock.MockCacheRepository,
		)
		input    getCategoryTestedInput
		expected getCategoryExpectedOutput
	}{
		{
			desc: "Success_FromCache",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(categorySerialized, nil)
			},
			input: getCategoryTestedInput{
				id: categoryID,
			},
			expected: getCategoryExpectedOutput{
				category: category,
				err:      nil,
			},
		},
		{
			desc: "Success_FromDB",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrInternal)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(categorySerialized), gomock.Eq(time.Duration(0))).
					Times(1).
					Return(nil)
			},
			input: getCategoryTestedInput{
				id: categoryID,
			},
			expected: getCategoryExpectedOutput{
				category: category,
				err:      nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrInternal)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
			},
			input: getCategoryTestedInput{
				id: categoryID,
			},
			expected: getCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrInternal)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: getCategoryTestedInput{
				id: categoryID,
			},
			expected: getCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrInternal)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(categorySerialized), gomock.Eq(time.Duration(0))).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: getCategoryTestedInput{
				id: categoryID,
			},
			expected: getCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_Deserialize",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return([]byte("invalid"), nil)
			},
			input: getCategoryTestedInput{
				id: categoryID,
			},
			expected: getCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			categoryRepo := mock.NewMockCategoryRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(categoryRepo, cache)

			categoryService := service.NewCategoryService(categoryRepo, cache)

			category, err := categoryService.GetCategory(ctx, tc.input.id)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.category, category, "Category mismatch")
		})
	}
}

type listCategoriesTestedInput struct {
	skip  uint64
	limit uint64
}

type listCategoriesExpectedOutput struct {
	categories []domain.Category
	err        error
}

func TestCategoryService_ListCategories(t *testing.T) {
	var categories []domain.Category

	for i := 0; i < 10; i++ {
		categories = append(categories, domain.Category{
			ID:   gofakeit.Uint64(),
			Name: gofakeit.ProductCategory(),
		})
	}

	ctx := context.Background()
	skip := gofakeit.Uint64()
	limit := gofakeit.Uint64()

	params := util.GenerateCacheKeyParams(skip, limit)
	cacheKey := util.GenerateCacheKey("categories", params)
	categoriesSerialized, _ := util.Serialize(categories)

	testCases := []struct {
		desc  string
		mocks func(
			categoryRepo *mock.MockCategoryRepository,
			cache *mock.MockCacheRepository,
		)
		input    listCategoriesTestedInput
		expected listCategoriesExpectedOutput
	}{
		{
			desc: "Success_FromCache",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(categoriesSerialized, nil)
			},
			input: listCategoriesTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listCategoriesExpectedOutput{
				categories: categories,
				err:        nil,
			},
		},
		{
			desc: "Success_FromDB",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrInternal)
				categoryRepo.EXPECT().
					ListCategories(gomock.Any(), gomock.Eq(skip), gomock.Eq(limit)).
					Times(1).
					Return(categories, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(categoriesSerialized), gomock.Eq(time.Duration(0))).
					Times(1).
					Return(nil)
			},
			input: listCategoriesTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listCategoriesExpectedOutput{
				categories: categories,
				err:        nil,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrInternal)
				categoryRepo.EXPECT().
					ListCategories(gomock.Any(), gomock.Eq(skip), gomock.Eq(limit)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: listCategoriesTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listCategoriesExpectedOutput{
				categories: nil,
				err:        domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrInternal)
				categoryRepo.EXPECT().
					ListCategories(gomock.Any(), gomock.Eq(skip), gomock.Eq(limit)).
					Times(1).
					Return(categories, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(categoriesSerialized), gomock.Eq(time.Duration(0))).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: listCategoriesTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listCategoriesExpectedOutput{
				categories: nil,
				err:        domain.ErrInternal,
			},
		},
		{
			desc: "Fail_Deserialize",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return([]byte("invalid"), nil)
			},
			input: listCategoriesTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listCategoriesExpectedOutput{
				categories: nil,
				err:        domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			categoryRepo := mock.NewMockCategoryRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(categoryRepo, cache)

			categoryService := service.NewCategoryService(categoryRepo, cache)

			categories, err := categoryService.ListCategories(ctx, tc.input.skip, tc.input.limit)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.categories, categories, "Categories mismatch")
		})
	}
}

type updateCategoryTestedInput struct {
	category *domain.Category
}

type updateCategoryExpectedOutput struct {
	category *domain.Category
	err      error
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	ctx := context.Background()
	categoryID := gofakeit.Uint64()
	categoryInput := &domain.Category{
		ID:   categoryID,
		Name: gofakeit.ProductCategory(),
	}
	categoryOutput := &domain.Category{
		ID:   categoryID,
		Name: categoryInput.Name,
	}
	existingCategory := &domain.Category{
		ID:   categoryID,
		Name: gofakeit.ProductCategory(),
	}

	cacheKey := util.GenerateCacheKey("category", categoryOutput.ID)
	categorySerialized, _ := util.Serialize(categoryOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			categoryRepo *mock.MockCategoryRepository,
			cache *mock.MockCacheRepository,
		)
		input    updateCategoryTestedInput
		expected updateCategoryExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryInput.ID)).
					Times(1).
					Return(existingCategory, nil)
				categoryRepo.EXPECT().
					UpdateCategory(gomock.Any(), gomock.Eq(categoryInput)).
					Times(1).
					Return(categoryOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(categorySerialized), gomock.Eq(ttl)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("categories:*")).
					Times(1).
					Return(nil)
			},
			input: updateCategoryTestedInput{
				category: categoryInput,
			},
			expected: updateCategoryExpectedOutput{
				category: categoryOutput,
				err:      nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryInput.ID)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
			},
			input: updateCategoryTestedInput{
				category: categoryInput,
			},
			expected: updateCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetByID",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryInput.ID)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: updateCategoryTestedInput{
				category: categoryInput,
			},
			expected: updateCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_EmptyData",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryInput.ID)).
					Times(1).
					Return(existingCategory, nil)
			},
			input: updateCategoryTestedInput{
				category: &domain.Category{
					ID: categoryInput.ID,
				},
			},
			expected: updateCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrNoUpdatedData,
			},
		},
		{
			desc: "Fail_SameData",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryInput.ID)).
					Times(1).
					Return(existingCategory, nil)
			},
			input: updateCategoryTestedInput{
				category: existingCategory,
			},
			expected: updateCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrNoUpdatedData,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryInput.ID)).
					Times(1).
					Return(existingCategory, nil)
				categoryRepo.EXPECT().
					UpdateCategory(gomock.Any(), gomock.Eq(categoryInput)).
					Times(1).
					Return(nil, domain.ErrConflictingData)
			},
			input: updateCategoryTestedInput{
				category: categoryInput,
			},
			expected: updateCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrConflictingData,
			},
		},
		{
			desc: "Fail_InternalErrorUpdate",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryInput.ID)).
					Times(1).
					Return(existingCategory, nil)
				categoryRepo.EXPECT().
					UpdateCategory(gomock.Any(), gomock.Eq(categoryInput)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: updateCategoryTestedInput{
				category: categoryInput,
			},
			expected: updateCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCache",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryInput.ID)).
					Times(1).
					Return(existingCategory, nil)
				categoryRepo.EXPECT().
					UpdateCategory(gomock.Any(), gomock.Eq(categoryInput)).
					Times(1).
					Return(categoryOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: updateCategoryTestedInput{
				category: categoryInput,
			},
			expected: updateCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryInput.ID)).
					Times(1).
					Return(existingCategory, nil)
				categoryRepo.EXPECT().
					UpdateCategory(gomock.Any(), gomock.Eq(categoryInput)).
					Times(1).
					Return(categoryOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(categorySerialized), gomock.Eq(ttl)).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: updateCategoryTestedInput{
				category: categoryInput,
			},
			expected: updateCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryInput.ID)).
					Times(1).
					Return(existingCategory, nil)
				categoryRepo.EXPECT().
					UpdateCategory(gomock.Any(), gomock.Eq(categoryInput)).
					Times(1).
					Return(categoryOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(categorySerialized), gomock.Eq(ttl)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("categories:*")).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: updateCategoryTestedInput{
				category: categoryInput,
			},
			expected: updateCategoryExpectedOutput{
				category: nil,
				err:      domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			categoryRepo := mock.NewMockCategoryRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(categoryRepo, cache)

			categoryService := service.NewCategoryService(categoryRepo, cache)

			category, err := categoryService.UpdateCategory(ctx, tc.input.category)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.category, category, "Category mismatch")
		})
	}
}

type deleteCategoryTestedInput struct {
	id uint64
}

type deleteCategoryExpectedOutput struct {
	err error
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	ctx := context.Background()
	categoryID := gofakeit.Uint64()

	cacheKey := util.GenerateCacheKey("category", categoryID)

	testCases := []struct {
		desc  string
		mocks func(
			categoryRepo *mock.MockCategoryRepository,
			cache *mock.MockCacheRepository,
		)
		input    deleteCategoryTestedInput
		expected deleteCategoryExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(&domain.Category{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("categories:*")).
					Times(1).
					Return(nil)
				categoryRepo.EXPECT().
					DeleteCategory(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(nil)
			},
			input: deleteCategoryTestedInput{
				id: categoryID,
			},
			expected: deleteCategoryExpectedOutput{
				err: nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
			},
			input: deleteCategoryTestedInput{
				id: categoryID,
			},
			expected: deleteCategoryExpectedOutput{
				err: domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetByID",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: deleteCategoryTestedInput{
				id: categoryID,
			},
			expected: deleteCategoryExpectedOutput{
				err: domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCache",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(&domain.Category{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: deleteCategoryTestedInput{
				id: categoryID,
			},
			expected: deleteCategoryExpectedOutput{
				err: domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(&domain.Category{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("categories:*")).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: deleteCategoryTestedInput{
				id: categoryID,
			},
			expected: deleteCategoryExpectedOutput{
				err: domain.ErrInternal,
			},
		},
		{
			desc: "Fail_InternalErrorDelete",
			mocks: func(
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(&domain.Category{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("categories:*")).
					Times(1).
					Return(nil)
				categoryRepo.EXPECT().
					DeleteCategory(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: deleteCategoryTestedInput{
				id: categoryID,
			},
			expected: deleteCategoryExpectedOutput{
				err: domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			categoryRepo := mock.NewMockCategoryRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(categoryRepo, cache)

			categoryService := service.NewCategoryService(categoryRepo, cache)

			err := categoryService.DeleteCategory(ctx, tc.input.id)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
		})
	}
}
