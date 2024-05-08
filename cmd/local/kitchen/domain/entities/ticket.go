package entities

import (
	"github.com/bmviniciuss/sagas-golang/pkg/utc"
	"github.com/google/uuid"
)

type TicketStatus string

func (os TicketStatus) String() string {
	return string(os)
}

const (
	TicketStatusCreatePending TicketStatus = "CREATE_PENDING"
)

type Ticket struct {
	ID           uuid.UUID
	CustomerID   uuid.UUID
	Amount       int64
	CurrencyCode string
	Status       TicketStatus
	CreatedAt    utc.Time
	UpdatedAt    utc.Time
	Items        []Item
}

type Item struct {
	ID        uuid.UUID
	Quantity  int32
	UnitPrice int64
}

func NewItem(id uuid.UUID, quantity int32, unitPrice int64) Item {
	return Item{
		ID:        id,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}
}

func NewTicket(customerID, globalID uuid.UUID, amount int64, currencyCode string, items []Item) *Ticket {
	return &Ticket{
		ID:           globalID,
		CustomerID:   customerID,
		Amount:       amount,
		CurrencyCode: currencyCode,
		Status:       TicketStatusCreatePending,
		CreatedAt:    utc.Now(),
		UpdatedAt:    utc.Now(),
		Items:        items,
	}
}
