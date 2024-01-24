package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

//go:generate mockgen -source=auth.go -destination=mock/auth.go -package=mock

// TokenService is an interface for interacting with token-related business logic
type TokenService interface {
	// CreateToken creates a new token for a given user
	CreateToken(user *domain.User) (string, error)
	// VerifyToken verifies the token and returns the payload
	VerifyToken(token string) (*domain.TokenPayload, error)
}

// UserService is an interface for interacting with user authentication-related business logic
type AuthService interface {
	// Login authenticates a user by email and password and returns a token
	Login(ctx context.Context, email, password string) (string, error)
}
