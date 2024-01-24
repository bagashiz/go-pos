package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

//go:generate mockgen -source=payment.go -destination=mock/payment.go -package=mock

// PaymentRepository is an interface for interacting with payment-related data
type PaymentRepository interface {
	// CreatePayment inserts a new payment into the database
	CreatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	// GetPaymentByID selects a payment by id
	GetPaymentByID(ctx context.Context, id uint64) (*domain.Payment, error)
	// ListPayments selects a list of payments with pagination
	ListPayments(ctx context.Context, skip, limit uint64) ([]domain.Payment, error)
	// UpdatePayment updates a payment
	UpdatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	// DeletePayment deletes a payment
	DeletePayment(ctx context.Context, id uint64) error
}

// PaymentService is an interface for interacting with payment-related business logic
type PaymentService interface {
	// CreatePayment creates a new payment
	CreatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	// GetPayment returns a payment by id
	GetPayment(ctx context.Context, id uint64) (*domain.Payment, error)
	// ListPayments returns a list of payments with pagination
	ListPayments(ctx context.Context, skip, limit uint64) ([]domain.Payment, error)
	// UpdatePayment updates a payment
	UpdatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	// DeletePayment deletes a payment
	DeletePayment(ctx context.Context, id uint64) error
}
