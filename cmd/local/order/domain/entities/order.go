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
	GlobalID     uuid.UUID
	CustomerID   uuid.UUID
	Amount       int64
	CurrencyCode string
	Status       OrderStatus
	Items        []Item
	CreatedAt    utc.Time
	UpdatedAt    utc.Time
}

type Item struct {
	ID        uuid.UUID
	Quantity  int16
	UnitPrice int64
}

func NewItem(id uuid.UUID, quantity int16, unitPrice int64) Item {
	return Item{
		ID:        id,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}
}

func NewOrder(customerID, globalID uuid.UUID, amount int64, currencyCode string, items []Item) Order {
	return Order{
		ID:           uuid.New(),
		GlobalID:     globalID,
		CustomerID:   customerID,
		Amount:       amount,
		CurrencyCode: currencyCode,
		Status:       OrderStatusApprovalPending,
		Items:        items,
		CreatedAt:    utc.Now(),
		UpdatedAt:    utc.Now(),
	}
}
