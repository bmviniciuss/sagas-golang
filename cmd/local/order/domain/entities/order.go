package entities

import (
	"github.com/bmviniciuss/sagas-golang/pkg/utc"
	"github.com/google/uuid"
)

type OrderStatus string

func (os OrderStatus) String() string {
	return string(os)
}

const (
	OrderStatusApprovalPending OrderStatus = "APPROVAL_PENDING"
)

type Order struct {
	ID           uuid.UUID
	CustomerID   uuid.UUID
	Amount       int64
	CurrencyCode string
	Status       OrderStatus
	CreatedAt    utc.Time
	UpdatedAt    utc.Time
}

func NewOrder(customerID, globalID uuid.UUID, amount int64, currencyCode string) Order {
	return Order{
		ID:           globalID,
		CustomerID:   customerID,
		Amount:       amount,
		CurrencyCode: currencyCode,
		Status:       OrderStatusApprovalPending,
		CreatedAt:    utc.Now(),
		UpdatedAt:    utc.Now(),
	}
}
