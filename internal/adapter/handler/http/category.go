package handler

import (
	"net/http"

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

// categoryResponse represents a category response body
type categoryResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// newCategoryResponse is a helper function to create a response body for handling category data
func newCategoryResponse(category *domain.Category) categoryResponse {
	return categoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}

// createCategoryRequest represents a request body for creating a new category
type createCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateCategory creates a new category
func (ch *CategoryHandler) CreateCategory(ctx *gin.Context) {
	var req createCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	category := domain.Category{
		Name: req.Name,
	}

	_, err := ch.svc.CreateCategory(ctx, &category)
	if err != nil {
		if err == domain.ErrConflictingData {
			errorResponse(ctx, http.StatusConflict, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newCategoryResponse(&category)

	successResponse(ctx, http.StatusCreated, rsp)
}

// getCategoryRequest represents a request body for retrieving a category
type getCategoryRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

// GetCategory retrieves a category by id
func (ch *CategoryHandler) GetCategory(ctx *gin.Context) {
	var req getCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	category, err := ch.svc.GetCategory(ctx, req.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newCategoryResponse(category)

	successResponse(ctx, http.StatusOK, rsp)
}

// listCategoriesRequest represents a request body for listing categories
type listCategoriesRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0"`
	Limit uint64 `form:"limit" binding:"required,min=5"`
}

// ListCategories lists all categories with pagination
func (ch *CategoryHandler) ListCategories(ctx *gin.Context) {
	var req listCategoriesRequest
	var categoriesList []categoryResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	categories, err := ch.svc.ListCategories(ctx, req.Skip, req.Limit)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	for _, category := range categories {
		categoriesList = append(categoriesList, newCategoryResponse(&category))
	}

	total := uint64(len(categoriesList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, categoriesList, "categories")

	successResponse(ctx, http.StatusOK, rsp)
}

// updateCategoryRequest represents a request body for updating a category
type updateCategoryRequest struct {
	Name string `json:"name" binding:"omitempty,required"`
}

// UpdateCategory updates a category
func (ch *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	var req updateCategoryRequest
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

	category := domain.Category{
		ID:   id,
		Name: req.Name,
	}

	_, err = ch.svc.UpdateCategory(ctx, &category)
	if err != nil {
		if err == domain.ErrDataNotFound {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		if err == domain.ErrNoUpdatedData {
			errorResponse(ctx, http.StatusBadRequest, err)
			return
		}

		if err == domain.ErrConflictingData {
			errorResponse(ctx, http.StatusConflict, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newCategoryResponse(&category)

	successResponse(ctx, http.StatusOK, rsp)
}

// deleteCategoryRequest represents a request body for deleting a category
type deleteCategoryRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

// DeleteCategory deletes a category
func (ch *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	var req deleteCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err := ch.svc.DeleteCategory(ctx, req.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	successResponse(ctx, http.StatusOK, nil)
}
