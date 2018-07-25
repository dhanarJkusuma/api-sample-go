package models

import (
	"database/sql/driver"
	"time"
)

type ExchangeAvgAggregate struct {
	ExchangeId int64    `json:"-"`
	Date       NullTime `json:"-"`
	DateString string   `json:"date"`
	Rate       float64  `json:"-"`
	RateString string   `json:"rate"`
	From       string   `json:"from"`
	To         string   `json:"to"`
	Average    float64  `json:"average"`
}

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
