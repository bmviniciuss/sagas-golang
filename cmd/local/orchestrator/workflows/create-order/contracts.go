package createorder

// Input is the data structure that represents the input data for the create order step
type Input struct {
	CustomerID   string      `json:"customer_id" validate:"required,uuid"`
	Card         string      `json:"card" validate:"required,uuid"`
	Amount       *int64      `json:"amount" validate:"required,gt=0"`
	CurrencyCode string      `json:"currency_code" validate:"required"`
	Items        []ItemInput `json:"items" validate:"required,min=1,dive"`
}

type ItemInput struct {
	ID        string `json:"id" validate:"required,uuid"`
	Quantity  *int32 `json:"quantity" validate:"required,gt=0"`
	UnitPrice *int64 `json:"unit_price" validate:"required,gt=0"`
}

// CreateOrderRequestPayload is the data structure that represents the payload for the create order step request
type CreateOrderRequestPayload struct {
	CustomerID   string                          `json:"customer_id" `
	Amount       *int64                          `json:"amount" `
	CurrencyCode string                          `json:"currency_code" `
	Items        []CreateOrderRequestItemPayload `json:"items" `
}

type CreateOrderRequestItemPayload struct {
	ID        string `json:"id" `
	Quantity  *int32 `json:"quantity" `
	UnitPrice *int64 `json:"unit_price" `
}

// CreateOrderResponsePayload is the data structure that represents the payload for the create order step response
type VerifyCustomerRequestPayload struct {
	CustomerID string `json:"customer_id" `
}

type AuthorizeCardRequestPayload struct {
	Card   string `json:"card"`
	Amount int64  `json:"amount"`
}
