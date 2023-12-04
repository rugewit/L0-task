package model

import "github.com/shopspring/decimal"

type Item struct {
	// chartID?
	ChrtID      int64           `json:"chrt_id" db:"chrt_id"`
	TrackNumber string          `json:"track_number" db:"track_number"`
	Price       decimal.Decimal `json:"price" db:"price"`
	RID         string          `json:"rid" db:"rid"`
	Name        string          `json:"name" db:"name"`
	Sale        decimal.Decimal `json:"sale" db:"sale"`
	Size        string          `json:"size" db:"size"`
	TotalPrice  decimal.Decimal `json:"total_price" db:"total_price"`
	NmID        int64           `json:"nm_id" db:"nm_id"`
	Brand       string          `json:"brand" db:"brand"`
	Status      int64           `json:"status" db:"status"`
}
