package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

//go:generate mockgen -source=order.go -destination=mock/order.go -package=mock

// OrderRepository is an interface for interacting with order-related data
type OrderRepository interface {
	// CreateOrder inserts a new order into the database
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	// GetOrderByID selects an order by id
	GetOrderByID(ctx context.Context, id uint64) (*domain.Order, error)
	// ListOrders selects a list of orders with pagination
	ListOrders(ctx context.Context, skip, limit uint64) ([]domain.Order, error)
}

// OrderService is an interface for interacting with order-related business logic
type OrderService interface {
	// CreateOrder creates a new order
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	// GetOrder returns an order by id
	GetOrder(ctx context.Context, id uint64) (*domain.Order, error)
	// ListOrders returns a list of orders with pagination
	ListOrders(ctx context.Context, skip, limit uint64) ([]domain.Order, error)
}
