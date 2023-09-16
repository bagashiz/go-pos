package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

// OrderRepository is an interface for interacting with order-related data
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetOrderByID(ctx context.Context, id uint64) (*domain.Order, error)
	ListOrders(ctx context.Context, skip, limit uint64) ([]domain.Order, error)
}

// OrderService is an interface for interacting with order-related business logic
type OrderService interface {
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetOrder(ctx context.Context, id uint64) (*domain.Order, error)
	ListOrders(ctx context.Context, skip, limit uint64) ([]domain.Order, error)
}
