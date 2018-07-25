package models

type ExchangeDataRequest struct {
	Id   int64   `json:"id"`
	Date string  `json:"date" validate:"required"`
	From string  `json:"from" validate:"required"`
	To   string  `json:"to" validate:"required"`
	Rate float64 `json:"rate" validate:"required"`
}
