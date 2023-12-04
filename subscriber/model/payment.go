package model

import "github.com/shopspring/decimal"

type Payment struct {
	Transaction  string          `json:"transaction" db:"transaction"`
	RequestID    string          `json:"request_id" db:"request_id"`
	Currency     string          `json:"currency" db:"currency"`
	Provider     string          `json:"provider" db:"provider"`
	Amount       int64           `json:"amount" db:"amount"`
	PaymentDT    int64           `json:"payment_dt" db:"payment_dt"`
	Bank         string          `json:"bank" db:"bank"`
	DeliveryCost decimal.Decimal `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   int64           `json:"goods_total" db:"goods_total"`
	CustomFee    decimal.Decimal `json:"custom_fee" db:"custom_fee"`
}
