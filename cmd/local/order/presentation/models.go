package presentation

import (
	"reflect"

	"github.com/bmviniciuss/sagas-golang/pkg/utc"
)

type Order struct {
	ID           string   `json:"id"`
	CustomerID   string   `json:"customer_id"`
	Amount       int64    `json:"amount"`
	CurrencyCode string   `json:"currency_code"`
	Status       string   `json:"status"`
	CreatedAt    utc.Time `json:"created_at"`
	UpdatedAt    utc.Time `json:"updated_at"`
}

type OrderList struct {
	Content []Order `json:"content"` // TODO: add pagination
}

type OrderById struct {
	ID           string   `json:"id"`
	CustomerID   string   `json:"customer_id"`
	Amount       int64    `json:"amount"`
	CurrencyCode string   `json:"currency_code"`
	Status       string   `json:"status"`
	CreatedAt    utc.Time `json:"created_at"`
	UpdatedAt    utc.Time `json:"updated_at"`
}

func (or *OrderById) IsEmpty() bool {
	return reflect.DeepEqual(OrderById{}, *or)
}
