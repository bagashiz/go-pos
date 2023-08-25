package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

// UserHandler represents the HTTP handler for user-related operations
type UserHandler struct {
	svc port.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(UserService port.UserService) *UserHandler {
	return &UserHandler{
		svc: UserService,
	}
}

// userResponse represents a user response body
type userResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// newUserResponse is a helper function to create a response body for handling user data
func newUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// registerRequest represents the request body for creating a user
type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Register creates a new user
func (uh *UserHandler) Register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	_, err := uh.svc.Register(ctx, &user)
	if err != nil {
		if err.Error() == "user already exists" {
			errorResponse(ctx, http.StatusConflict, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newUserResponse(&user)

	successResponse(ctx, http.StatusCreated, rsp)
}

// listUsersRequest represents the request body for listing users
type listUsersRequest struct {
	PageID   uint64 `form:"page_id" binding:"required,min=1"`
	PageSize uint64 `form:"page_size" binding:"required,min=5"`
}

// ListUsers lists all users with pagination
func (uh *UserHandler) ListUsers(ctx *gin.Context) {
	var req listUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	users, err := uh.svc.ListUsers(ctx, req.PageID, req.PageSize)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	data := make([]userResponse, 0)
	for _, user := range users {
		data = append(data, newUserResponse(user))
	}

	meta := newMeta(uint64(len(data)), req.PageSize, req.PageID)

	rsp := make(map[string]any)
	rsp["meta"] = meta
	rsp["users"] = data

	successResponse(ctx, http.StatusOK, rsp)
}

// getUserRequest represents the request body for getting a user
type getUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

// GetUser gets a user by ID
func (uh *UserHandler) GetUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := uh.svc.GetUser(ctx, req.ID)
	if err != nil {
		if err.Error() == "user not found" {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newUserResponse(user)

	successResponse(ctx, http.StatusOK, rsp)
}

// updateUserRequest represents the request body for updating a user
type updateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// UpdateUser updates a user's name, email, and password
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user := domain.User{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	_, err = uh.svc.UpdateUser(ctx, &user)
	if err != nil {
		if err.Error() == "user not found" {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newUserResponse(&user)

	successResponse(ctx, http.StatusOK, rsp)
}

// deleteUserRequest represents the request body for deleting a user
type deleteUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

// DeleteUser deletes a user by ID
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err := uh.svc.DeleteUser(ctx, req.ID)
	if err != nil {
		if err.Error() == "user not found" {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	successResponse(ctx, http.StatusOK, nil)
}
