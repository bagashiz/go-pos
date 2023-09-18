package port

import (
	"context"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

// TokenService is an interface for interacting with token-related business logic
type TokenService interface {
	CreateToken(user *domain.User, duration time.Duration) (string, error)
	VerifyToken(token string) (*domain.TokenPayload, error)
}

// UserService is an interface for interacting with user authentication-related business logic
type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
}
