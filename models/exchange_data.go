package models

import "time"

type ExchangeData struct {
	Id         int64     `json:"id"`
	ExchangeId int64     `json:"exchangeRateId"`
	Date       time.Time `json:"date"`
	Rate       float64   `json:"rate"`
}
