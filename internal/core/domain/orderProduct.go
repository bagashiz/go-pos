package domain

import "time"

// OrderProduct is an entity that represents pivot table between order and product
type OrderProduct struct {
	ID         uint64    `json:"id"`
	OrderID    uint64    `json:"order_id"`
	ProductID  uint64    `json:"product_id"`
	Quantity   int64     `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Order      *Order    `json:"order"`
	Product    *Product  `json:"product"`
}
