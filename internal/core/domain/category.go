package domain

import "time"

// Category is an entity that represents a category of product
type Category struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
