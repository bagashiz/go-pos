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
	ID        uint64      `json:"id"`
	Name      string      `json:"name"`
	Type      PaymentType `json:"type"`
	Logo      string      `json:"logo"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
