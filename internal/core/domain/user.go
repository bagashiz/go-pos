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
	ID        uint64
	Name      string
	Email     string
	Password  string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}
