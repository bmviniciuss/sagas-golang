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
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusFailed    OrderStatus = "failed"
)

type Order struct {
	ID           uuid.UUID
	GlobalID     uuid.UUID
	ClientID     uuid.UUID
	CustomerID   uuid.UUID
	Total        int64
	CurrencyCode string
	Status       OrderStatus
	CreatedAt    utc.Time
	UpdatedAt    utc.Time
}

func NewOrder(clientID, customerID, globalID uuid.UUID, total int64, currencyCode string) Order {
	return Order{
		ID:           uuid.New(),
		GlobalID:     globalID,
		ClientID:     clientID,
		CustomerID:   customerID,
		Total:        total,
		CurrencyCode: currencyCode,
		Status:       OrderStatusPending,
		CreatedAt:    utc.Now(),
		UpdatedAt:    utc.Now(),
	}
}
