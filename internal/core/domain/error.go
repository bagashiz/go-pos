package domain

import (
	"errors"
	"strings"
)

var (
	// ErrDataNotFound is an error for when requested data is not found
	ErrDataNotFound = errors.New("data not found")
	// ErrNoUpdatedData is an error for when no data is provided to update
	ErrNoUpdatedData = errors.New("no data to update")
	// ErrConflictingData is an error for when data conflicts with existing data
	ErrConflictingData = errors.New("data conflicts with existing data in unique column")
	// ErrInsufficientStock is an error for when product stock is not enough
	ErrInsufficientStock = errors.New("product stock is not enough")
	// ErrInsufficientPayment is an error for when total paid is less than total price
	ErrInsufficientPayment = errors.New("total paid is less than total price")
)

// IsUniqueConstraintViolationError checks if the error is a unique constraint violation error
func IsUniqueConstraintViolationError(err error) bool {
	return strings.Contains(err.Error(), "23505")
}
