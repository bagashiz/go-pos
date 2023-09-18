package service

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
)

/**
 * ProductService implements port.ProductService and port.CategoryService
 * interfaces and provides an access to the product and category repositories
 */
type ProductService struct {
	productRepo  port.ProductRepository
	categoryRepo port.CategoryRepository
}

// NewProductService creates a new product service instance
func NewProductService(productRepo port.ProductRepository, categoryRepo port.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo,
		categoryRepo,
	}
}

// CreateProduct creates a new product
func (ps *ProductService) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	category, err := ps.categoryRepo.GetCategoryByID(ctx, product.CategoryID)
	if err != nil {
		return nil, err
	}

	product.Category = category

	_, err = ps.productRepo.CreateProduct(ctx, product)
	if err != nil {
		if port.IsUniqueConstraintViolationError(err) {
			return nil, port.ErrConflictingData
		}
	}

	return product, nil
}

// GetProduct retrieves a product by id
func (ps *ProductService) GetProduct(ctx context.Context, id uint64) (*domain.Product, error) {
	product, err := ps.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category, err := ps.categoryRepo.GetCategoryByID(ctx, product.CategoryID)
	if err != nil {
		return nil, err
	}

	product.Category = category

	return product, nil
}

// ListProducts retrieves a list of products
func (ps *ProductService) ListProducts(ctx context.Context, search string, categoryId, skip, limit uint64) ([]domain.Product, error) {
	products, err := ps.productRepo.ListProducts(ctx, search, categoryId, skip, limit)
	if err != nil {
		return nil, err
	}

	for i, product := range products {
		category, err := ps.categoryRepo.GetCategoryByID(ctx, product.CategoryID)
		if err != nil {
			return nil, err
		}

		products[i].Category = category
	}

	return products, nil
}

// UpdateProduct updates a product
func (ps *ProductService) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	existingProduct, err := ps.productRepo.GetProductByID(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	emptyData := product.CategoryID == 0 && product.Name == "" && product.Image == "" && product.Price == 0 && product.Stock == 0
	sameData := existingProduct.CategoryID == product.CategoryID && existingProduct.Name == product.Name && existingProduct.Image == product.Image && existingProduct.Price == product.Price && existingProduct.Stock == product.Stock
	if emptyData || sameData {
		return nil, port.ErrNoUpdatedData
	}

	if product.CategoryID == 0 {
		product.CategoryID = existingProduct.CategoryID
	}

	category, err := ps.categoryRepo.GetCategoryByID(ctx, product.CategoryID)
	if err != nil {
		return nil, err
	}

	product.Category = category

	_, err = ps.productRepo.UpdateProduct(ctx, product)
	if err != nil {
		if port.IsUniqueConstraintViolationError(err) {
			return nil, port.ErrConflictingData
		}
	}

	return product, nil
}

// DeleteProduct deletes a product
func (ps *ProductService) DeleteProduct(ctx context.Context, id uint64) error {
	_, err := ps.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return err
	}

	return ps.productRepo.DeleteProduct(ctx, id)
}
