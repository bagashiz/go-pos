package handler

import (
	"net/http"

	"github.com/bagashiz/go-pos/internal/core/domain"
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
	AccessToken string `json:"token"`
}

// newAuthResponse is a helper function to create a response body for handling authentication data
func newAuthResponse(token string) authResponse {
	return authResponse{
		AccessToken: token,
	}
}

// loginRequest represents the request body for logging in a user
type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Login gives a registered user an access token if the credentials are valid
func (ah *AuthHandler) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	token, err := ah.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			errorResponse(ctx, http.StatusUnauthorized, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newAuthResponse(token)

	successResponse(ctx, http.StatusOK, rsp)
}
