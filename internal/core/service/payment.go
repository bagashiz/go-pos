package service

import (
	"context"
	"errors"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
)

/**
 * PaymentService implements port.PaymentService interface
 * and provides an access to the payment repository
 */
type PaymentService struct {
	repo port.PaymentRepository
}

// NewPaymentService creates a new payment service instance
func NewPaymentService(repo port.PaymentRepository) *PaymentService {
	return &PaymentService{
		repo,
	}
}

// CreatePayment creates a new payment
func (ps *PaymentService) CreatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	return ps.repo.CreatePayment(ctx, payment)
}

// GetPayment retrieves a payment by id
func (ps *PaymentService) GetPayment(ctx context.Context, id uint64) (*domain.Payment, error) {
	return ps.repo.GetPaymentByID(ctx, id)
}

// ListPayments retrieves a list of payments
func (ps *PaymentService) ListPayments(ctx context.Context, skip, limit uint64) ([]*domain.Payment, error) {
	return ps.repo.ListPayments(ctx, skip, limit)
}

// UpdatePayment updates a payment
func (ps *PaymentService) UpdatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	_, err := ps.repo.GetPaymentByID(ctx, payment.ID)
	if err != nil {
		return nil, err
	}

	emptyData := payment.Name == "" && payment.Type == "" && payment.Logo == ""
	if emptyData {
		return nil, errors.New("no data to update")
	}

	return ps.repo.UpdatePayment(ctx, payment)
}

// DeletePayment deletes a payment
func (ps *PaymentService) DeletePayment(ctx context.Context, id uint64) error {
	_, err := ps.repo.GetPaymentByID(ctx, id)
	if err != nil {
		return err
	}

	return ps.repo.DeletePayment(ctx, id)
}
