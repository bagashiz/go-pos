package handler

import (
	"net/http"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

// UserHandler represents the HTTP handler for user-related requests
type UserHandler struct {
	svc port.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(svc port.UserService) *UserHandler {
	return &UserHandler{
		svc,
	}
}

// userResponse represents a user response body
type userResponse struct {
	ID        uint64    `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"test@example.com"`
	CreatedAt time.Time `json:"c0eated_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
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
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	create a new user account with default role "cashier"
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			registerRequest	body		registerRequest	true	"Register request"
//	@Success		201				{object}	userResponse	"User created"
//	@Failure		400				{object}	response		"Validation error"
//	@Failure		401				{object}	response		"Unauthorized error"
//	@Failure		404				{object}	response		"Data not found error"
//	@Failure		409				{object}	response		"Data conflict error"
//	@Failure		500				{object}	response		"Internal server error"
//	@Router			/users [post]
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
		handleError(ctx, err)
		return
	}

	rsp := newUserResponse(&user)

	successResponse(ctx, http.StatusCreated, rsp)
}

// listUsersRequest represents the request body for listing users
type listUsersRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListUsers godoc
//
//	@Summary		List users
//	@Description	List users with pagination
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64		true	"Skip"
//	@Param			limit	query		uint64		true	"Limit"
//	@Success		200		{object}	response	"Users displayed"
//	@Failure		400		{object}	response	"Validation error"
//	@Failure		500		{object}	response	"Internal server error"
//	@Router			/users [get]
func (uh *UserHandler) ListUsers(ctx *gin.Context) {
	var req listUsersRequest
	var usersList []userResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	users, err := uh.svc.ListUsers(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, user := range users {
		usersList = append(usersList, newUserResponse(&user))
	}

	total := uint64(len(usersList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, usersList, "users")

	successResponse(ctx, http.StatusOK, rsp)
}

// getUserRequest represents the request body for getting a user
type getUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetUser godoc
//
//	@Summary		Get a user
//	@Description	Get a user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"User ID"
//	@Success		200	{object}	userResponse	"User displayed"
//	@Failure		400	{object}	response		"Validation error"
//	@Failure		404	{object}	response		"Data not found error"
//	@Failure		500	{object}	response		"Internal server error"
//	@Router			/users/{id} [get]
func (uh *UserHandler) GetUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := uh.svc.GetUser(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newUserResponse(user)

	successResponse(ctx, http.StatusOK, rsp)
}

// updateUserRequest represents the request body for updating a user
type updateUserRequest struct {
	Name     string          `json:"name" binding:"omitempty,required" example:"John Doe"`
	Email    string          `json:"email" binding:"omitempty,required,email" example:"test@example.com"`
	Password string          `json:"password" binding:"omitempty,required,min=8" example:"12345678"`
	Role     domain.UserRole `json:"role" binding:"omitempty,required,user_role" example:"admin"`
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Update a user's name, email, password, or role by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id					path		uint64				true	"User ID"
//	@Param			updateUserRequest	body		updateUserRequest	true	"Update user request"
//	@Success		200					{object}	userResponse		"User updated"
//	@Failure		400					{object}	response			"Validation error"
//	@Failure		401					{object}	response			"Unauthorized error"
//	@Failure		404					{object}	response			"Data not found error"
//	@Failure		500					{object}	response			"Internal server error"
//	@Router			/users/{id} [put]
//	@Security		BearerAuth
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var req updateUserRequest
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

	user := domain.User{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	_, err = uh.svc.UpdateUser(ctx, &user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newUserResponse(&user)

	successResponse(ctx, http.StatusOK, rsp)
}

// deleteUserRequest represents the request body for deleting a user
type deleteUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64		true	"User ID"
//	@Success		200	{object}	response	"User deleted"
//	@Failure		400	{object}	response	"Validation error"
//	@Failure		401	{object}	response	"Unauthorized error"
//	@Failure		404	{object}	response	"Data not found error"
//	@Failure		500	{object}	response	"Internal server error"
//	@Router			/users/{id} [delete]
//	@Security		BearerAuth
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err := uh.svc.DeleteUser(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	successResponse(ctx, http.StatusOK, nil)
}
