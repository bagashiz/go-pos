package domain

import (
	"time"
)

// UserRole is an enum for user's role
type UserRole string

// UserRole enum values
const (
	Admin   UserRole = "admin"
	Cashier UserRole = "cashier"
)

// User is an entity that represents a user
type User struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
