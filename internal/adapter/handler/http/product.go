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
	ID        uint64           `json:"id" example:"1"`
	SKU       string           `json:"sku" example:"9a4c25d3-9786-492c-b084-85cb75c1ee3e"`
	Name      string           `json:"name" example:"Chiki Ball"`
	Stock     int64            `json:"stock" example:"100"`
	Price     float64          `json:"price" example:"5000"`
	Image     string           `json:"image" example:"https://example.com/chiki-ball.png"`
	Category  categoryResponse `json:"category"`
	CreatedAt time.Time        `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time        `json:"updated_at" example:"1970-01-01T00:00:00Z"`
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
	CategoryID uint64  `json:"category_id" binding:"required,min=1" example:"1"`
	Name       string  `json:"name" binding:"required" example:"Chiki Ball"`
	Image      string  `json:"image" binding:"required" example:"https://example.com/chiki-ball.png"`
	Price      float64 `json:"price" binding:"required,min=0" example:"5000"`
	Stock      int64   `json:"stock" binding:"required,min=0" example:"100"`
}

// CreateProduct godoc
//
//	@Summary		Create a new product
//	@Description	create a new product with name, image, price, and stock
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			createProductRequest	body		createProductRequest	true	"Create product request"
//	@Success		201						{object}	productResponse			"Product created"
//	@Failure		400						{object}	response				"Validation error"
//	@Failure		401						{object}	response				"Unauthorized error"
//	@Failure		404						{object}	response				"Data not found error"
//	@Failure		409						{object}	response				"Data conflict error"
//	@Failure		500						{object}	response				"Internal server error"
//	@Router			/products [post]
//	@Security		BearerAuth
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
		handleError(ctx, err)
		return
	}

	rsp := newProductResponse(&product)

	successResponse(ctx, http.StatusCreated, rsp)
}

// getProductRequest represents a request body for retrieving a product
type getProductRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetProduct godoc
//
//	@Summary		Get a product
//	@Description	get a product by id with its category
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Product ID"
//	@Success		200	{object}	productResponse	"Product retrieved"
//	@Failure		400	{object}	response		"Validation error"
//	@Failure		404	{object}	response		"Data not found error"
//	@Failure		500	{object}	response		"Internal server error"
//	@Router			/products/{id} [get]
func (ph *ProductHandler) GetProduct(ctx *gin.Context) {
	var req getProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	product, err := ph.svc.GetProduct(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newProductResponse(product)

	successResponse(ctx, http.StatusOK, rsp)
}

// listProductsRequest represents a request body for listing products
type listProductsRequest struct {
	CategoryID uint64 `form:"category_id" binding:"omitempty,min=1" example:"1"`
	Query      string `form:"q" binding:"omitempty" example:"Chiki"`
	Skip       uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit      uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListProducts godoc
//
//	@Summary		List products
//	@Description	List products with pagination
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			category_id	query		uint64		false	"Category ID"
//	@Param			q			query		string		false	"Query"
//	@Param			skip		query		uint64		true	"Skip"
//	@Param			limit		query		uint64		true	"Limit"
//	@Success		200			{object}	response	"Products retrieved"
//	@Failure		400			{object}	response	"Validation error"
//	@Failure		500			{object}	response	"Internal server error"
//	@Router			/products [get]
func (ph *ProductHandler) ListProducts(ctx *gin.Context) {
	var req listProductsRequest
	var productsList []productResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	products, err := ph.svc.ListProducts(ctx, req.Query, req.CategoryID, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
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
	CategoryID uint64  `json:"category_id" binding:"omitempty,required,min=1" example:"1"`
	Name       string  `json:"name" binding:"omitempty,required" example:"Nutrisari Jeruk"`
	Image      string  `json:"image" binding:"omitempty,required" example:"https://example.com/nutrisari-jeruk.png"`
	Price      float64 `json:"price" binding:"omitempty,required,min=0" example:"2000"`
	Stock      int64   `json:"stock" binding:"omitempty,required,min=0" example:"200"`
}

// UpdateProduct godoc
//
//	@Summary		Update a product
//	@Description	update a product's name, image, price, or stock by id
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id						path		uint64					true	"Product ID"
//	@Param			updateProductRequest	body		updateProductRequest	true	"Update product request"
//	@Success		200						{object}	productResponse			"Product updated"
//	@Failure		400						{object}	response				"Validation error"
//	@Failure		401						{object}	response				"Unauthorized error"
//	@Failure		404						{object}	response				"Data not found error"
//	@Failure		409						{object}	response				"Data conflict error"
//	@Failure		500						{object}	response				"Internal server error"
//	@Router			/products/{id} [put]
//	@Security		BearerAuth
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
		handleError(ctx, err)
		return
	}

	rsp := newProductResponse(&product)

	successResponse(ctx, http.StatusOK, rsp)
}

// deleteProductRequest represents a request body for deleting a product
type deleteProductRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// DeleteProduct godoc
//
//	@Summary		Delete a product
//	@Description	Delete a product by id
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64		true	"Product ID"
//	@Success		200	{object}	response	"Product deleted"
//	@Failure		400	{object}	response	"Validation error"
//	@Failure		401	{object}	response	"Unauthorized error"
//	@Failure		404	{object}	response	"Data not found error"
//	@Failure		500	{object}	response	"Internal server error"
//	@Router			/products/{id} [delete]
//	@Security		BearerAuth
func (ph *ProductHandler) DeleteProduct(ctx *gin.Context) {
	var req deleteProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err := ph.svc.DeleteProduct(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	successResponse(ctx, http.StatusOK, nil)
}
