package service

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/bagashiz/go-pos/internal/core/util"
)

/**
 * PaymentService implements port.PaymentService interface
 * and provides an access to the payment repository
 * and cache service
 */
type PaymentService struct {
	repo  port.PaymentRepository
	cache port.CacheRepository
}

// NewPaymentService creates a new payment service instance
func NewPaymentService(repo port.PaymentRepository, cache port.CacheRepository) *PaymentService {
	return &PaymentService{
		repo,
		cache,
	}
}

// CreatePayment creates a new payment
func (ps *PaymentService) CreatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	_, err := ps.repo.CreatePayment(ctx, payment)
	if err != nil {
		if port.IsUniqueConstraintViolationError(err) {
			return nil, port.ErrConflictingData
		}

		return nil, err
	}

	cacheKey := util.GenerateCacheKey("payment", payment.ID)
	paymentSerialized, err := util.Serialize(payment)
	if err != nil {
		return nil, err
	}

	err = ps.cache.Set(ctx, cacheKey, paymentSerialized, 0)
	if err != nil {
		return nil, err
	}

	err = ps.cache.DeleteByPrefix(ctx, "payments:*")
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetPayment retrieves a payment by id
func (ps *PaymentService) GetPayment(ctx context.Context, id uint64) (*domain.Payment, error) {
	var payment *domain.Payment

	cacheKey := util.GenerateCacheKey("payment", id)
	cachedPayment, err := ps.cache.Get(ctx, cacheKey)
	if err == nil {
		err = util.Deserialize(cachedPayment, &payment)
		if err != nil {
			return nil, err
		}

		return payment, nil
	}

	payment, err = ps.repo.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	paymentSerialized, err := util.Serialize(payment)
	if err != nil {
		return nil, err
	}

	err = ps.cache.Set(ctx, cacheKey, paymentSerialized, 0)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// ListPayments retrieves a list of payments
func (ps *PaymentService) ListPayments(ctx context.Context, skip, limit uint64) ([]domain.Payment, error) {
	var payments []domain.Payment

	params := util.GenerateCacheKeyParams(skip, limit)
	cacheKey := util.GenerateCacheKey("payments", params)

	cachedPayments, err := ps.cache.Get(ctx, cacheKey)
	if err == nil {
		err = util.Deserialize(cachedPayments, &payments)
		if err != nil {
			return nil, err
		}

		return payments, nil
	}

	payments, err = ps.repo.ListPayments(ctx, skip, limit)
	if err != nil {
		return nil, err
	}

	paymentsSerialized, err := util.Serialize(payments)
	if err != nil {
		return nil, err
	}

	err = ps.cache.Set(ctx, cacheKey, paymentsSerialized, 0)
	if err != nil {
		return nil, err
	}

	return payments, nil

}

// UpdatePayment updates a payment
func (ps *PaymentService) UpdatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	existingPayment, err := ps.repo.GetPaymentByID(ctx, payment.ID)
	if err != nil {
		return nil, err
	}

	emptyData := payment.Name == "" && payment.Type == "" && payment.Logo == ""
	sameData := existingPayment.Name == payment.Name && existingPayment.Type == payment.Type && existingPayment.Logo == payment.Logo
	if emptyData || sameData {
		return nil, port.ErrNoUpdatedData
	}

	_, err = ps.repo.UpdatePayment(ctx, payment)
	if err != nil {
		if port.IsUniqueConstraintViolationError(err) {
			return nil, port.ErrConflictingData
		}

		return nil, err
	}

	cacheKey := util.GenerateCacheKey("payment", payment.ID)
	paymentSerialized, err := util.Serialize(payment)
	if err != nil {
		return nil, err
	}

	err = ps.cache.Set(ctx, cacheKey, paymentSerialized, 0)
	if err != nil {
		return nil, err
	}

	err = ps.cache.DeleteByPrefix(ctx, "payments:*")
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// DeletePayment deletes a payment
func (ps *PaymentService) DeletePayment(ctx context.Context, id uint64) error {
	_, err := ps.repo.GetPaymentByID(ctx, id)
	if err != nil {
		return err
	}

	cacheKey := util.GenerateCacheKey("payment", id)
	_ = ps.cache.Delete(ctx, cacheKey)

	err = ps.cache.DeleteByPrefix(ctx, "payments:*")
	if err != nil {
		return err
	}

	return ps.repo.DeletePayment(ctx, id)
}
