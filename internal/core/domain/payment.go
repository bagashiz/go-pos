package domain

import "time"

// PaymentType is an enum for payment's type
type PaymentType string

// PaymentType enum values
const (
	Cash    PaymentType = "CASH"
	EWallet PaymentType = "E-WALLET"
	EDC     PaymentType = "EDC"
)

// Payment is an entity that represents a payment
type Payment struct {
	ID        uint64
	Name      string
	Type      PaymentType
	Logo      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
