package handler

import (
	"net/http"

	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

// AuthHandler represents the HTTP handler for authentication-related requests
type AuthHandler struct {
	svc port.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(svc port.AuthService) *AuthHandler {
	return &AuthHandler{
		svc,
	}
}

// authResponse represents an authentication response body
type authResponse struct {
	AccessToken string `json:"token" example:"v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."`
}

// newAuthResponse is a helper function to create a response body for handling authentication data
func newAuthResponse(token string) authResponse {
	return authResponse{
		AccessToken: token,
	}
}

// loginRequest represents the request body for logging in a user
type loginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678" minLength:"8"`
}

// Login godoc
//
//	@Summary		Login and get an access token
//	@Description	Logs in a registered user and returns an access token if the credentials are valid.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest	true	"Login request body"
//	@Success		200		{object}	authResponse	"Succesfully logged in"
//	@Failure		400		{object}	response		"Validation error"
//	@Failure		401		{object}	response		"Unauthorized error"
//	@Failure		500		{object}	response		"Internal server error"
//	@Router			/users/login [post]
func (ah *AuthHandler) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	token, err := ah.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newAuthResponse(token)

	successResponse(ctx, http.StatusOK, rsp)
}
