package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

// PaymentRepository is an interface for interacting with payment-related data
type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	GetPaymentByID(ctx context.Context, id uint64) (*domain.Payment, error)
	ListPayments(ctx context.Context, skip, limit uint64) ([]*domain.Payment, error)
	UpdatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	DeletePayment(ctx context.Context, id uint64) error
}

// PaymentService is an interface for interacting with payment-related business logic
type PaymentService interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	GetPayment(ctx context.Context, id uint64) (*domain.Payment, error)
	ListPayments(ctx context.Context, skip, limit uint64) ([]*domain.Payment, error)
	UpdatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	DeletePayment(ctx context.Context, id uint64) error
}
