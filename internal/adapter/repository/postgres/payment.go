package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/core/domain"
)

/**
 * PaymentRepository implements port.PaymentRepository interface
 * and provides an access to the postgres database
 */
type PaymentRepository struct {
	db *DB
}

// NewPaymentRepository creates a new payment repository instance
func NewPaymentRepository(db *DB) *PaymentRepository {
	return &PaymentRepository{
		db,
	}
}

// CreatePayment creates a new payment record in the database
func (pr *PaymentRepository) CreatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	query := psql.Insert("payments").
		Columns("name", "type", "logo").
		Values(payment.Name, payment.Type, payment.Logo).
		Suffix("RETURNING *")

	err := query.QueryRowContext(ctx).Scan(
		&payment.ID,
		&payment.Name,
		&payment.Type,
		&payment.Logo,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetPaymentByID retrieves a payment record from the database by id
func (pr *PaymentRepository) GetPaymentByID(ctx context.Context, id uint64) (*domain.Payment, error) {
	query := psql.Select("*").
		From("payments").
		Where(sq.Eq{"id": id}).
		Limit(1)

	var payment domain.Payment

	err := query.QueryRowContext(ctx).Scan(
		&payment.ID,
		&payment.Name,
		&payment.Type,
		&payment.Logo,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	return &payment, nil
}

// ListPayments retrieves a list of payments from the database
func (pr *PaymentRepository) ListPayments(ctx context.Context, skip, limit uint64) ([]*domain.Payment, error) {
	query := psql.Select("*").
		From("payments").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var payments []*domain.Payment

	for rows.Next() {
		var payment domain.Payment

		err := rows.Scan(
			&payment.ID,
			&payment.Name,
			&payment.Type,
			&payment.Logo,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, &payment)
	}

	return payments, nil
}

// UpdatePayment updates a payment record in the database
func (pr *PaymentRepository) UpdatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	name := nullString(payment.Name)
	paymentType := nullString(payment.Type)
	logo := nullString(payment.Logo)

	query := psql.Update("payments").
		Set("name", sq.Expr("COALESCE(?, name)", name)).
		Set("type", sq.Expr("COALESCE(?, type)", paymentType)).
		Set("logo", sq.Expr("COALESCE(?, logo)", logo)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": payment.ID}).
		Suffix("RETURNING *")

	err := query.QueryRowContext(ctx).Scan(
		&payment.ID,
		&payment.Name,
		&payment.Type,
		&payment.Logo,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// DeletePayment deletes a payment record from the database by id
func (pr *PaymentRepository) DeletePayment(ctx context.Context, id uint64) error {
	query := psql.Delete("payments").
		Where(sq.Eq{"id": id})

	_, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
