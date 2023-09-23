package domain

import (
	"time"

	"github.com/google/uuid"
)

// Product is an entity that represents a product
type Product struct {
	ID         uint64    `json:"id"`
	CategoryID uint64    `json:"category_id"`
	SKU        uuid.UUID `json:"sku"`
	Name       string    `json:"name"`
	Stock      int64     `json:"stock"`
	Price      float64   `json:"price"`
	Image      string    `json:"image"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Category   *Category `json:"category"`
}
