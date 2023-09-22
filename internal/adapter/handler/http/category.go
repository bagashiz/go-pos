package handler

import (
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

// CategoryHandler represents the HTTP handler for category-related requests
type CategoryHandler struct {
	svc port.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewCategoryHandler(svc port.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		svc,
	}
}

// createCategoryRequest represents a request body for creating a new category
type createCategoryRequest struct {
	Name string `json:"name" binding:"required" example:"Foods"`
}

// CreateCategory godoc
//
//	@Summary		Create a new category
//	@Description	create a new category with name
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			createCategoryRequest	body		createCategoryRequest	true	"Create category request"
//	@Success		200						{object}	categoryResponse		"Category created"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		401						{object}	errorResponse			"Unauthorized error"
//	@Failure		404						{object}	errorResponse			"Data not found error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/categories [post]
//	@Security		BearerAuth
func (ch *CategoryHandler) CreateCategory(ctx *gin.Context) {
	var req createCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	category := domain.Category{
		Name: req.Name,
	}

	_, err := ch.svc.CreateCategory(ctx, &category)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newCategoryResponse(&category)

	handleSuccess(ctx, rsp)
}

// getCategoryRequest represents a request body for retrieving a category
type getCategoryRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetCategory godoc
//
//	@Summary		Get a category
//	@Description	get a category by id
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64				true	"Category ID"
//	@Success		200	{object}	categoryResponse	"Category retrieved"
//	@Failure		400	{object}	errorResponse		"Validation error"
//	@Failure		404	{object}	errorResponse		"Data not found error"
//	@Failure		500	{object}	errorResponse		"Internal server error"
//	@Router			/categories/{id} [get]
//	@Security		BearerAuth
func (ch *CategoryHandler) GetCategory(ctx *gin.Context) {
	var req getCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	category, err := ch.svc.GetCategory(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newCategoryResponse(category)

	handleSuccess(ctx, rsp)
}

// listCategoriesRequest represents a request body for listing categories
type listCategoriesRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListCategories godoc
//
//	@Summary		List categories
//	@Description	List categories with pagination
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			true	"Skip"
//	@Param			limit	query		uint64			true	"Limit"
//	@Success		200		{object}	meta			"Categories displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/categories [get]
//	@Security		BearerAuth
func (ch *CategoryHandler) ListCategories(ctx *gin.Context) {
	var req listCategoriesRequest
	var categoriesList []categoryResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	categories, err := ch.svc.ListCategories(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, category := range categories {
		categoriesList = append(categoriesList, newCategoryResponse(&category))
	}

	total := uint64(len(categoriesList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, categoriesList, "categories")

	handleSuccess(ctx, rsp)
}

// updateCategoryRequest represents a request body for updating a category
type updateCategoryRequest struct {
	Name string `json:"name" binding:"omitempty,required" example:"Beverages"`
}

// UpdateCategory godoc
//
//	@Summary		Update a category
//	@Description	update a category's name by id
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id						path		uint64					true	"Category ID"
//	@Param			updateCategoryRequest	body		updateCategoryRequest	true	"Update category request"
//	@Success		200						{object}	categoryResponse		"Category updated"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		401						{object}	errorResponse			"Unauthorized error"
//	@Failure		404						{object}	errorResponse			"Data not found error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/categories/{id} [put]
//	@Security		BearerAuth
func (ch *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	var req updateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := stringToUint64(idStr)
	if err != nil {
		validationError(ctx, err)
		return
	}

	category := domain.Category{
		ID:   id,
		Name: req.Name,
	}

	_, err = ch.svc.UpdateCategory(ctx, &category)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newCategoryResponse(&category)

	handleSuccess(ctx, rsp)
}

// deleteCategoryRequest represents a request body for deleting a category
type deleteCategoryRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// DeleteCategory godoc
//
//	@Summary		Delete a category
//	@Description	Delete a category by id
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Category ID"
//	@Success		200	{object}	response		"Category deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/categories/{id} [delete]
//	@Security		BearerAuth
func (ch *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	var req deleteCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := ch.svc.DeleteCategory(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
