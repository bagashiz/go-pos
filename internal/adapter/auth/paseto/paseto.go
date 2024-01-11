package paseto

import (
	"errors"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

/**
 * PasetoToken implements port.TokenService interface
 * and provides an access to the paseto library
 */
type PasetoToken struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// New creates a new paseto instance
func New(symmetricKey string) (port.TokenService, error) {
	validSymmetricKey := len(symmetricKey) == chacha20poly1305.KeySize
	if !validSymmetricKey {
		return nil, errors.New("invalid token key size")
	}

	return &PasetoToken{
		paseto.NewV2(),
		[]byte(symmetricKey),
	}, nil
}

// CreateToken creates a new paseto token
func (pt *PasetoToken) CreateToken(user *domain.User, duration time.Duration) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	payload := domain.TokenPayload{
		ID:        id,
		UserID:    user.ID,
		Role:      user.Role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	token, err := pt.paseto.Encrypt(pt.symmetricKey, payload, nil)

	return token, err

}

// VerifyToken verifies the paseto token
func (pt *PasetoToken) VerifyToken(token string) (*domain.TokenPayload, error) {
	var payload domain.TokenPayload

	err := pt.paseto.Decrypt(token, pt.symmetricKey, &payload, nil)
	if err != nil {
		return nil, port.ErrInvalidToken
	}

	isExpired := time.Now().After(payload.ExpiredAt)
	if isExpired {
		return nil, port.ErrExpiredToken
	}

	return &payload, nil
}
