package createorder

// first step
type Input struct {
	CustomerID   string `mapstructure:"customer_id"`
	Date         string `mapstructure:"date"`
	Total        *int64 `mapstructure:"total"`
	CurrencyCode string `mapstructure:"currency_code"`
}
