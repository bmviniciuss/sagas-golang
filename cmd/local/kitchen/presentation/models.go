package presentation

import (
	"reflect"

	"github.com/bmviniciuss/sagas-golang/pkg/utc"
)

type Ticket struct {
	ID           string   `json:"id"`
	CustomerID   string   `json:"customer_id"`
	Amount       int64    `json:"amount"`
	CurrencyCode string   `json:"currency_code"`
	Status       string   `json:"status"`
	CreatedAt    utc.Time `json:"created_at"`
	UpdatedAt    utc.Time `json:"updated_at"`
}

type TicketsList struct {
	Content []Ticket `json:"content"` // TODO: add pagination
}

type TicketByID struct {
	ID           string       `json:"id"`
	CustomerID   string       `json:"customer_id"`
	Amount       int64        `json:"amount"`
	CurrencyCode string       `json:"currency_code"`
	Status       string       `json:"status"`
	Items        []TicketItem `json:"items"`
	CreatedAt    utc.Time     `json:"created_at"`
	UpdatedAt    utc.Time     `json:"updated_at"`
}

func (or *TicketByID) IsEmpty() bool {
	return reflect.DeepEqual(TicketByID{}, *or)
}

type TicketItem struct {
	ID        string `json:"id"`
	Quantity  int32  `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
}
