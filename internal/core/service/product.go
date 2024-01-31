package service

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/bagashiz/go-pos/internal/core/util"
)

/**
 * ProductService implements port.ProductService and port.CategoryService
 * interfaces and provides an access to the product and category repositories
 * and cache service
 */
type ProductService struct {
	productRepo  port.ProductRepository
	categoryRepo port.CategoryRepository
	cache        port.CacheRepository
}

// NewProductService creates a new product service instance
func NewProductService(productRepo port.ProductRepository, categoryRepo port.CategoryRepository, cache port.CacheRepository) *ProductService {
	return &ProductService{
		productRepo,
		categoryRepo,
		cache,
	}
}

// CreateProduct creates a new product
func (ps *ProductService) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	category, err := ps.categoryRepo.GetCategoryByID(ctx, product.CategoryID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	product.Category = category

	product, err = ps.productRepo.CreateProduct(ctx, product)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	cacheKey := util.GenerateCacheKey("product", product.ID)
	productSerialized, err := util.Serialize(product)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = ps.cache.Set(ctx, cacheKey, productSerialized, 0)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = ps.cache.DeleteByPrefix(ctx, "products:*")
	if err != nil {
		return nil, domain.ErrInternal
	}

	return product, nil
}

// GetProduct retrieves a product by id
func (ps *ProductService) GetProduct(ctx context.Context, id uint64) (*domain.Product, error) {
	var product *domain.Product

	cacheKey := util.GenerateCacheKey("product", id)
	cachedProduct, err := ps.cache.Get(ctx, cacheKey)
	if err == nil {
		err := util.Deserialize(cachedProduct, &product)
		if err != nil {
			return nil, domain.ErrInternal
		}
		return product, nil
	}

	product, err = ps.productRepo.GetProductByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	category, err := ps.categoryRepo.GetCategoryByID(ctx, product.CategoryID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	product.Category = category

	productSerialized, err := util.Serialize(product)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = ps.cache.Set(ctx, cacheKey, productSerialized, 0)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return product, nil
}

// ListProducts retrieves a list of products
func (ps *ProductService) ListProducts(ctx context.Context, search string, categoryID, skip, limit uint64) ([]domain.Product, error) {
	var products []domain.Product

	params := util.GenerateCacheKeyParams(skip, limit, categoryID, search)
	cacheKey := util.GenerateCacheKey("products", params)

	cachedProducts, err := ps.cache.Get(ctx, cacheKey)
	if err == nil {
		err := util.Deserialize(cachedProducts, &products)
		if err != nil {
			return nil, domain.ErrInternal
		}
		return products, nil
	}

	products, err = ps.productRepo.ListProducts(ctx, search, categoryID, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	for i, product := range products {
		category, err := ps.categoryRepo.GetCategoryByID(ctx, product.CategoryID)
		if err != nil {
			if err == domain.ErrDataNotFound {
				return nil, err
			}
			return nil, domain.ErrInternal
		}

		products[i].Category = category
	}

	productsSerialized, err := util.Serialize(products)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = ps.cache.Set(ctx, cacheKey, productsSerialized, 0)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return products, nil
}

// UpdateProduct updates a product
func (ps *ProductService) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	existingProduct, err := ps.productRepo.GetProductByID(ctx, product.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	emptyData := product.CategoryID == 0 &&
		product.Name == "" &&
		product.Image == "" &&
		product.Price == 0 &&
		product.Stock == 0

	sameData := existingProduct.CategoryID == product.CategoryID &&
		existingProduct.Name == product.Name &&
		existingProduct.Image == product.Image &&
		existingProduct.Price == product.Price &&
		existingProduct.Stock == product.Stock

	if emptyData || sameData {
		return nil, domain.ErrNoUpdatedData
	}

	if product.CategoryID == 0 {
		product.CategoryID = existingProduct.CategoryID
	}

	category, err := ps.categoryRepo.GetCategoryByID(ctx, product.CategoryID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	product.Category = category

	_, err = ps.productRepo.UpdateProduct(ctx, product)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	cacheKey := util.GenerateCacheKey("product", product.ID)

	err = ps.cache.Delete(ctx, cacheKey)
	if err != nil {
		return nil, domain.ErrInternal
	}

	productSerialized, err := util.Serialize(product)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = ps.cache.Set(ctx, cacheKey, productSerialized, 0)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = ps.cache.DeleteByPrefix(ctx, "products:*")
	if err != nil {
		return nil, domain.ErrInternal
	}

	return product, nil
}

// DeleteProduct deletes a product
func (ps *ProductService) DeleteProduct(ctx context.Context, id uint64) error {
	_, err := ps.productRepo.GetProductByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	cacheKey := util.GenerateCacheKey("product", id)

	err = ps.cache.Delete(ctx, cacheKey)
	if err != nil {
		return domain.ErrInternal
	}

	err = ps.cache.DeleteByPrefix(ctx, "products:*")
	if err != nil {
		return domain.ErrInternal
	}

	return ps.productRepo.DeleteProduct(ctx, id)
}
