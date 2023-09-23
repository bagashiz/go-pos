package domain

import (
	"time"

	"github.com/google/uuid"
)

// TokenPayload is an entity that represents the payload of the token
type TokenPayload struct {
	ID        uuid.UUID `json:"id"`
	UserID    uint64    `json:"user_id"`
	Role      UserRole  `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
