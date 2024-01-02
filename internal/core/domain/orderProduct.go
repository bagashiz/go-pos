package domain

import "time"

// OrderProduct is an entity that represents pivot table between order and product
type OrderProduct struct {
	ID         uint64
	OrderID    uint64
	ProductID  uint64
	Quantity   int64
	TotalPrice float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Order      *Order
	Product    *Product
}
