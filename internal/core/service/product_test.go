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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type createProductTestedInput struct {
	product *domain.Product
}

type createProductExpectedOutput struct {
	product *domain.Product
	err     error
}

func TestProductService_CreateProduct(t *testing.T) {
	ctx := context.Background()
	categoryID := gofakeit.Uint64()
	categoryName := gofakeit.ProductCategory()
	category := &domain.Category{
		ID:   categoryID,
		Name: categoryName,
	}

	productName := gofakeit.ProductName()
	productStock := gofakeit.Int64()
	productPrice := gofakeit.Float64()
	productImage := gofakeit.ImageURL(400, 400)
	productSKU, _ := uuid.NewUUID()

	productInput := &domain.Product{
		Name:       productName,
		Stock:      productStock,
		Price:      productPrice,
		Image:      productImage,
		CategoryID: categoryID,
	}

	productOutput := &domain.Product{
		ID:         gofakeit.Uint64(),
		SKU:        productSKU,
		Name:       productName,
		Stock:      productStock,
		Price:      productPrice,
		Image:      productImage,
		CategoryID: categoryID,
		Category:   category,
		CreatedAt:  gofakeit.Date(),
		UpdatedAt:  gofakeit.Date(),
	}

	cacheKey := util.GenerateCacheKey("product", productOutput.ID)
	productSerialized, _ := util.Serialize(productOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			productRepo *mock.MockProductRepository,
			categoryRepo *mock.MockCategoryRepository,
			cache *mock.MockCacheRepository,
		)
		input    createProductTestedInput
		expected createProductExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				productRepo.EXPECT().
					CreateProduct(gomock.Any(), gomock.Eq(productInput)).
					Times(1).
					Return(productOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(productSerialized), gomock.Eq(ttl)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("products:*")).
					Times(1).
					Return(nil)
			},
			input: createProductTestedInput{
				product: productInput,
			},
			expected: createProductExpectedOutput{
				product: productOutput,
				err:     nil,
			},
		},
		{
			desc: "Fail_NotFoudGetCategory",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
			},
			input: createProductTestedInput{
				product: productInput,
			},
			expected: createProductExpectedOutput{
				product: nil,
				err:     domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetCategory",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: createProductTestedInput{
				product: productInput,
			},
			expected: createProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				productRepo.EXPECT().
					CreateProduct(gomock.Any(), gomock.Eq(productInput)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: createProductTestedInput{
				product: productInput,
			},
			expected: createProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				productRepo.EXPECT().
					CreateProduct(gomock.Any(), gomock.Eq(productInput)).
					Times(1).
					Return(nil, domain.ErrConflictingData)
			},
			input: createProductTestedInput{
				product: productInput,
			},
			expected: createProductExpectedOutput{
				product: nil,
				err:     domain.ErrConflictingData,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				productRepo.EXPECT().
					CreateProduct(gomock.Any(), gomock.Eq(productInput)).
					Times(1).
					Return(productOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(productSerialized), gomock.Eq(ttl)).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: createProductTestedInput{
				product: productInput,
			},
			expected: createProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				productRepo.EXPECT().
					CreateProduct(gomock.Any(), gomock.Eq(productInput)).
					Times(1).
					Return(productOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(productSerialized), gomock.Eq(ttl)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("products:*")).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: createProductTestedInput{
				product: productInput,
			},
			expected: createProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		// tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			// TODO: fix race condition to enable parallel testing
			// t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productRepo := mock.NewMockProductRepository(ctrl)
			categoryRepo := mock.NewMockCategoryRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(productRepo, categoryRepo, cache)

			productService := service.NewProductService(productRepo, categoryRepo, cache)

			product, err := productService.CreateProduct(ctx, tc.input.product)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.product, product, "Product mismatch")
		})
	}
}

type getProductTestedInput struct {
	id uint64
}

type getProductExpectedOutput struct {
	product *domain.Product
	err     error
}

func TestProductService_GetProduct(t *testing.T) {
	ctx := context.Background()
	productID := gofakeit.Uint64()
	productSKU, _ := uuid.NewUUID()
	categoryID := gofakeit.Uint64()
	categoryName := gofakeit.ProductCategory()
	category := &domain.Category{
		ID:   categoryID,
		Name: categoryName,
	}

	productOutput := &domain.Product{
		ID:         productID,
		SKU:        productSKU,
		Name:       gofakeit.ProductName(),
		Stock:      gofakeit.Int64(),
		Price:      gofakeit.Float64(),
		Image:      gofakeit.ImageURL(400, 400),
		CategoryID: categoryID,
		Category:   category,
		CreatedAt:  gofakeit.Date(),
		UpdatedAt:  gofakeit.Date(),
	}

	cacheKey := util.GenerateCacheKey("product", productOutput.ID)
	productSerialized, _ := util.Serialize(productOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			productRepo *mock.MockProductRepository,
			categoryRepo *mock.MockCategoryRepository,
			cache *mock.MockCacheRepository,
		)
		input    getProductTestedInput
		expected getProductExpectedOutput
	}{
		{
			desc: "Success_FromCache",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(productSerialized, nil)
			},
			input: getProductTestedInput{
				id: productID,
			},
			expected: getProductExpectedOutput{
				product: productOutput,
				err:     nil,
			},
		},
		{
			desc: "Success_FromDB",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(productOutput, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(productOutput.CategoryID)).
					Times(1).
					Return(category, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(productSerialized), gomock.Eq(ttl)).
					Times(1).
					Return(nil)
			},
			input: getProductTestedInput{
				id: productID,
			},
			expected: getProductExpectedOutput{
				product: productOutput,
				err:     nil,
			},
		},
		{
			desc: "Fail_Deserialize",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return([]byte("invalid"), nil)
			},
			input: getProductTestedInput{
				id: productID,
			},
			expected: getProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrInternal)
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: getProductTestedInput{
				id: productID,
			},
			expected: getProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
			},
			input: getProductTestedInput{
				id: productID,
			},
			expected: getProductExpectedOutput{
				product: nil,
				err:     domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetCategory",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(productOutput, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(productOutput.CategoryID)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: getProductTestedInput{
				id: productID,
			},
			expected: getProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_NotFoundGetCategory",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(productOutput, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(productOutput.CategoryID)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
			},
			input: getProductTestedInput{
				id: productID,
			},
			expected: getProductExpectedOutput{
				product: nil,
				err:     domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(productOutput, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(productOutput.CategoryID)).
					Times(1).
					Return(category, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(productSerialized), gomock.Eq(ttl)).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: getProductTestedInput{
				id: productID,
			},
			expected: getProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		// tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			// TODO: fix race condition to enable parallel testing
			// t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productRepo := mock.NewMockProductRepository(ctrl)
			categoryRepo := mock.NewMockCategoryRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(productRepo, categoryRepo, cache)

			productService := service.NewProductService(productRepo, categoryRepo, cache)

			product, err := productService.GetProduct(ctx, tc.input.id)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.product, product, "Product mismatch")
		})
	}
}

type listProductsTestedInput struct {
	search     string
	categoryID uint64
	skip       uint64
	limit      uint64
}

type listProductsExpectedOutput struct {
	products []domain.Product
	err      error
}

func TestProductService_ListProducts(t *testing.T) {
	var products []domain.Product

	categoryID := gofakeit.Uint64()
	categoryName := gofakeit.ProductCategory()
	category := &domain.Category{
		ID:   categoryID,
		Name: categoryName,
	}

	for i := 0; i < 10; i++ {
		productSKU, _ := uuid.NewUUID()
		products = append(products, domain.Product{
			ID:         gofakeit.Uint64(),
			SKU:        productSKU,
			Name:       gofakeit.ProductName(),
			Stock:      gofakeit.Int64(),
			Price:      gofakeit.Float64(),
			Image:      gofakeit.ImageURL(400, 400),
			CategoryID: categoryID,
			Category:   category,
			CreatedAt:  gofakeit.Date(),
			UpdatedAt:  gofakeit.Date(),
		})
	}

	ctx := context.Background()
	skip := gofakeit.Uint64()
	limit := gofakeit.Uint64()
	search := ""

	params := util.GenerateCacheKeyParams(skip, limit, categoryID, search)
	cacheKey := util.GenerateCacheKey("products", params)
	productsSerialized, _ := util.Serialize(products)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			productRepo *mock.MockProductRepository,
			categoryRepo *mock.MockCategoryRepository,
			cache *mock.MockCacheRepository,
		)
		input    listProductsTestedInput
		expected listProductsExpectedOutput
	}{
		{
			desc: "Success_FromCache",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(productsSerialized, nil)
			},
			input: listProductsTestedInput{
				search:     search,
				categoryID: categoryID,
				skip:       skip,
				limit:      limit,
			},
			expected: listProductsExpectedOutput{
				products: products,
				err:      nil,
			},
		},
		{
			desc: "Success_FromDB",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
				productRepo.EXPECT().
					ListProducts(gomock.Any(), gomock.Eq(search), gomock.Eq(categoryID), gomock.Eq(skip), gomock.Eq(limit)).
					Times(1).
					Return(products, nil)
				for i := range products {
					categoryRepo.EXPECT().
						GetCategoryByID(gomock.Any(), gomock.Eq(products[i].CategoryID)).
						Times(1).
						Return(category, nil)
				}
				productsSerialized, _ := util.Serialize(products)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(productsSerialized), gomock.Eq(ttl)).
					Times(1).
					Return(nil)
			},
			input: listProductsTestedInput{
				search:     search,
				categoryID: categoryID,
				skip:       skip,
				limit:      limit,
			},
			expected: listProductsExpectedOutput{
				products: products,
				err:      nil,
			},
		},
		{
			desc: "Fail_Deserialize",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return([]byte("invalid"), nil)
			},
			input: listProductsTestedInput{
				search:     search,
				categoryID: categoryID,
				skip:       skip,
				limit:      limit,
			},
			expected: listProductsExpectedOutput{
				products: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
				productRepo.EXPECT().
					ListProducts(gomock.Any(), gomock.Eq(search), gomock.Eq(categoryID), gomock.Eq(skip), gomock.Eq(limit)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: listProductsTestedInput{
				search:     search,
				categoryID: categoryID,
				skip:       skip,
				limit:      limit,
			},
			expected: listProductsExpectedOutput{
				products: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_NotFoundGetCategory",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
				productRepo.EXPECT().
					ListProducts(gomock.Any(), gomock.Eq(search), gomock.Eq(categoryID), gomock.Eq(skip), gomock.Eq(limit)).
					Times(1).
					Return(products, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(products[0].CategoryID)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
			},
			input: listProductsTestedInput{
				search:     search,
				categoryID: categoryID,
				skip:       skip,
				limit:      limit,
			},
			expected: listProductsExpectedOutput{
				products: nil,
				err:      domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetCategory",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
				productRepo.EXPECT().
					ListProducts(gomock.Any(), gomock.Eq(search), gomock.Eq(categoryID), gomock.Eq(skip), gomock.Eq(limit)).
					Times(1).
					Return(products, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(products[0].CategoryID)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: listProductsTestedInput{
				search:     search,
				categoryID: categoryID,
				skip:       skip,
				limit:      limit,
			},
			expected: listProductsExpectedOutput{
				products: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
				productRepo.EXPECT().
					ListProducts(gomock.Any(), gomock.Eq(search), gomock.Eq(categoryID), gomock.Eq(skip), gomock.Eq(limit)).
					Times(1).
					Return(products, nil)
				for i := range products {
					categoryRepo.EXPECT().
						GetCategoryByID(gomock.Any(), gomock.Eq(products[i].CategoryID)).
						Times(1).
						Return(category, nil)
				}
				productsSerialized, _ := util.Serialize(products)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(productsSerialized), gomock.Eq(ttl)).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: listProductsTestedInput{
				search:     search,
				categoryID: categoryID,
				skip:       skip,
				limit:      limit,
			},
			expected: listProductsExpectedOutput{
				products: nil,
				err:      domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		// tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			// TODO: fix race condition to enable parallel testing
			// t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productRepo := mock.NewMockProductRepository(ctrl)
			categoryRepo := mock.NewMockCategoryRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(productRepo, categoryRepo, cache)

			productService := service.NewProductService(productRepo, categoryRepo, cache)

			products, err := productService.ListProducts(ctx, tc.input.search, tc.input.categoryID, tc.input.skip, tc.input.limit)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.products, products, "Products mismatch")
		})
	}
}

type updateProductTestedInput struct {
	product *domain.Product
}

type updateProductExpectedOutput struct {
	product *domain.Product
	err     error
}

func TestProductService_UpdateProduct(t *testing.T) {
	ctx := context.Background()
	productID := gofakeit.Uint64()
	productSKU, _ := uuid.NewUUID()
	categoryID := gofakeit.Uint64()
	categoryName := gofakeit.ProductCategory()
	category := &domain.Category{
		ID:   categoryID,
		Name: categoryName,
	}

	productName := gofakeit.ProductName()
	productStock := gofakeit.Int64()
	productPrice := gofakeit.Float64()
	productImage := gofakeit.ImageURL(400, 400)

	productInput := &domain.Product{
		ID:         productID,
		SKU:        productSKU,
		Name:       productName,
		Stock:      productStock,
		Price:      productPrice,
		Image:      productImage,
		CategoryID: categoryID,
	}

	productOutput := &domain.Product{
		ID:         productID,
		SKU:        productSKU,
		Name:       productName,
		Stock:      productStock,
		Price:      productPrice,
		Image:      productImage,
		CategoryID: categoryID,
		Category:   category,
	}

	existingProduct := &domain.Product{
		ID:    productID,
		SKU:   productSKU,
		Name:  gofakeit.ProductName(),
		Stock: gofakeit.Int64(),
		Price: gofakeit.Float64(),
		Image: gofakeit.ImageURL(400, 400),
	}

	cacheKey := util.GenerateCacheKey("product", productOutput.ID)
	productSerialized, _ := util.Serialize(productOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			productRepo *mock.MockProductRepository,
			categoryRepo *mock.MockCategoryRepository,
			cache *mock.MockCacheRepository,
		)
		input    updateProductTestedInput
		expected updateProductExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(existingProduct, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				productRepo.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Eq(productInput)).
					Times(1).
					Return(productOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(productSerialized), gomock.Eq(ttl)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("products:*")).
					Times(1).
					Return(nil)
			},
			input: updateProductTestedInput{
				product: productInput,
			},
			expected: updateProductExpectedOutput{
				product: productOutput,
				err:     nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
			},
			input: updateProductTestedInput{
				product: productInput,
			},
			expected: updateProductExpectedOutput{
				product: nil,
				err:     domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: updateProductTestedInput{
				product: productInput,
			},
			expected: updateProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_EmptyData",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(existingProduct, nil)
			},
			input: updateProductTestedInput{
				product: &domain.Product{
					ID: productID,
				},
			},
			expected: updateProductExpectedOutput{
				product: nil,
				err:     domain.ErrNoUpdatedData,
			},
		},
		{
			desc: "Fail_SameData",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(existingProduct, nil)
			},
			input: updateProductTestedInput{
				product: existingProduct,
			},
			expected: updateProductExpectedOutput{
				product: nil,
				err:     domain.ErrNoUpdatedData,
			},
		},
		{
			desc: "Fail_NotFoundGetCategory",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(existingProduct, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
			},
			input: updateProductTestedInput{
				product: productInput,
			},
			expected: updateProductExpectedOutput{
				product: nil,
				err:     domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetCategory",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(existingProduct, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: updateProductTestedInput{
				product: productInput,
			},
			expected: updateProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(existingProduct, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				productRepo.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Eq(productInput)).
					Times(1).
					Return(nil, domain.ErrConflictingData)
			},
			input: updateProductTestedInput{
				product: productInput,
			},
			expected: updateProductExpectedOutput{
				product: nil,
				err:     domain.ErrConflictingData,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(existingProduct, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				productRepo.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Eq(productInput)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: updateProductTestedInput{
				product: productInput,
			},
			expected: updateProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCache",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(existingProduct, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				productRepo.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Eq(productInput)).
					Times(1).
					Return(productOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: updateProductTestedInput{
				product: productInput,
			},
			expected: updateProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(existingProduct, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				productRepo.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Eq(productInput)).
					Times(1).
					Return(productOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(productSerialized), gomock.Eq(ttl)).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: updateProductTestedInput{
				product: productInput,
			},
			expected: updateProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {

				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(existingProduct, nil)
				categoryRepo.EXPECT().
					GetCategoryByID(gomock.Any(), gomock.Eq(categoryID)).
					Times(1).
					Return(category, nil)
				productRepo.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Eq(productInput)).
					Times(1).
					Return(productOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(productSerialized), gomock.Eq(ttl)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("products:*")).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: updateProductTestedInput{
				product: productInput,
			},
			expected: updateProductExpectedOutput{
				product: nil,
				err:     domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		// tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			// TODO: fix race condition to enable parallel testing
			// t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productRepo := mock.NewMockProductRepository(ctrl)
			categoryRepo := mock.NewMockCategoryRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(productRepo, categoryRepo, cache)

			productService := service.NewProductService(productRepo, categoryRepo, cache)

			product, err := productService.UpdateProduct(ctx, tc.input.product)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.product, product, "Product mismatch")
		})
	}
}

type deleteProductTestedInput struct {
	id uint64
}

type deleteProductExpectedOutput struct {
	err error
}

func TestProductService_DeleteProduct(t *testing.T) {
	ctx := context.Background()
	productID := gofakeit.Uint64()

	cacheKey := util.GenerateCacheKey("product", productID)

	testCases := []struct {
		desc  string
		mocks func(
			productRepo *mock.MockProductRepository,
			categoryRepo *mock.MockCategoryRepository,
			cache *mock.MockCacheRepository,
		)
		input    deleteProductTestedInput
		expected deleteProductExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(nil, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("products:*")).
					Times(1).
					Return(nil)
				productRepo.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(nil)
			},
			input: deleteProductTestedInput{
				id: productID,
			},
			expected: deleteProductExpectedOutput{
				err: nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(nil, domain.ErrDataNotFound)
			},
			input: deleteProductTestedInput{
				id: productID,
			},
			expected: deleteProductExpectedOutput{
				err: domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetByID",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(nil, domain.ErrInternal)
			},
			input: deleteProductTestedInput{
				id: productID,
			},
			expected: deleteProductExpectedOutput{
				err: domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCache",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(nil, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: deleteProductTestedInput{
				id: productID,
			},
			expected: deleteProductExpectedOutput{
				err: domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(nil, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("products:*")).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: deleteProductTestedInput{
				id: productID,
			},
			expected: deleteProductExpectedOutput{
				err: domain.ErrInternal,
			},
		},
		{
			desc: "Fail_InternalErrorDelete",
			mocks: func(
				productRepo *mock.MockProductRepository,
				categoryRepo *mock.MockCategoryRepository,
				cache *mock.MockCacheRepository,
			) {
				productRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(nil, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Times(1).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("products:*")).
					Times(1).
					Return(nil)
				productRepo.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Eq(productID)).
					Times(1).
					Return(domain.ErrInternal)
			},
			input: deleteProductTestedInput{
				id: productID,
			},
			expected: deleteProductExpectedOutput{
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

			productRepo := mock.NewMockProductRepository(ctrl)
			categoryRepo := mock.NewMockCategoryRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(productRepo, categoryRepo, cache)

			productService := service.NewProductService(productRepo, categoryRepo, cache)

			err := productService.DeleteProduct(ctx, tc.input.id)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
		})
	}
}
