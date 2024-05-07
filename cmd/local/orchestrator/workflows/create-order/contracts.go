package createorder

type Input struct {
	CustomerID   string `json:"customer_id" validate:"required,uuid"`
	Amount       *int64 `json:"amount" validate:"required,gt=0"`
	CurrencyCode string `json:"currency_code" validate:"required"`
	Items        []Item `json:"items" validate:"required,min=1,dive"`
}

type Item struct {
	ID        string `json:"id" validate:"required,uuid"`
	Quantity  *int16 `json:"quantity" validate:"required,gt=0"`
	UnitPrice *int64 `json:"unit_price" validate:"required,gt=0"`
}

type CreateOrderRequestPayload struct {
	CustomerID   string                          `json:"customer_id" `
	Amount       *int64                          `json:"amount" `
	CurrencyCode string                          `json:"currency_code" `
	Items        []CreateOrderRequestItemPayload `json:"items" `
}

type CreateOrderRequestItemPayload struct {
	ID        string `json:"id" `
	Quantity  *int16 `json:"quantity" `
	UnitPrice *int64 `json:"unit_price" `
}
