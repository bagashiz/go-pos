package handler

import (
	"net/http"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

// ProductHandler represents the HTTP handler for product-related requests
type ProductHandler struct {
	svc port.ProductService
}

// NewProductHandler creates a new ProductHandler instance
func NewProductHandler(svc port.ProductService) *ProductHandler {
	return &ProductHandler{
		svc,
	}
}

// productResponse represents a product response body
type productResponse struct {
	ID        uint64           `json:"id"`
	SKU       string           `json:"sku"`
	Name      string           `json:"name"`
	Stock     int64            `json:"stock"`
	Price     float64          `json:"price"`
	Image     string           `json:"image"`
	Category  categoryResponse `json:"category"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// newProductResponse is a helper function to create a response body for handling product data
func newProductResponse(product *domain.Product) productResponse {
	return productResponse{
		ID:        product.ID,
		SKU:       product.SKU.String(),
		Name:      product.Name,
		Stock:     product.Stock,
		Price:     product.Price,
		Image:     product.Image,
		Category:  newCategoryResponse(product.Category),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

// createProductRequest represents a request body for creating a new product
type createProductRequest struct {
	CategoryID uint64  `json:"category_id" binding:"required,min=1"`
	Name       string  `json:"name" binding:"required"`
	Image      string  `json:"image" binding:"required"`
	Price      float64 `json:"price" binding:"required,min=0"`
	Stock      int64   `json:"stock" binding:"required,min=0"`
}

// CreateProduct creates a new product
func (ph *ProductHandler) CreateProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	product := domain.Product{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Image:      req.Image,
		Price:      req.Price,
		Stock:      req.Stock,
	}

	_, err := ph.svc.CreateProduct(ctx, &product)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newProductResponse(&product)

	successResponse(ctx, http.StatusCreated, rsp)
}

// getProductRequest represents a request body for retrieving a product
type getProductRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

// GetProduct retrieves a product by id
func (ph *ProductHandler) GetProduct(ctx *gin.Context) {
	var req getProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	product, err := ph.svc.GetProduct(ctx, req.ID)
	if err != nil {
		if err.Error() == "product not found" {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newProductResponse(product)

	successResponse(ctx, http.StatusOK, rsp)
}

// listProductsRequest represents a request body for listing products
type listProductsRequest struct {
	CategoryID uint64 `form:"category_id" binding:"omitempty,min=1"`
	Query      string `form:"q" binding:"omitempty"`
	Skip       uint64 `form:"skip" binding:"required,min=0"`
	Limit      uint64 `form:"limit" binding:"required,min=5"`
}

// ListProducts lists all products with pagination
func (ph *ProductHandler) ListProducts(ctx *gin.Context) {
	var req listProductsRequest
	var productsList []productResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	products, err := ph.svc.ListProducts(ctx, req.Query, req.CategoryID, req.Skip, req.Limit)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	for _, product := range products {
		productsList = append(productsList, newProductResponse(&product))
	}

	total := uint64(len(productsList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, productsList, "products")

	successResponse(ctx, http.StatusOK, rsp)
}

// updateProductRequest represents a request body for updating a product
type updateProductRequest struct {
	CategoryID uint64  `json:"category_id" binding:"omitempty,required,min=1"`
	Name       string  `json:"name" binding:"omitempty,required"`
	Image      string  `json:"image" binding:"omitempty,required"`
	Price      float64 `json:"price" binding:"omitempty,required,min=0"`
	Stock      int64   `json:"stock" binding:"omitempty,required,min=0"`
}

// UpdateProduct updates a product
func (ph *ProductHandler) UpdateProduct(ctx *gin.Context) {
	var req updateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := stringToUint64(idStr)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	product := domain.Product{
		ID:         id,
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Image:      req.Image,
		Price:      req.Price,
		Stock:      req.Stock,
	}

	_, err = ph.svc.UpdateProduct(ctx, &product)
	if err != nil {
		if err.Error() == "product not found" {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		if err.Error() == "no data to update" {
			errorResponse(ctx, http.StatusBadRequest, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newProductResponse(&product)

	successResponse(ctx, http.StatusOK, rsp)
}

// deleteProductRequest represents a request body for deleting a product
type deleteProductRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

// DeleteProduct deletes a product
func (ph *ProductHandler) DeleteProduct(ctx *gin.Context) {
	var req deleteProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err := ph.svc.DeleteProduct(ctx, req.ID)
	if err != nil {
		if err.Error() == "product not found" {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	successResponse(ctx, http.StatusOK, nil)
}
