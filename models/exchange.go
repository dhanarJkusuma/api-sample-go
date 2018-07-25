package models

type Exchange struct {
	Id   int64  `json:"id"`
	From string `json:"from" validate:"required"`
	To   string `json:"to" validate:"required"`
}
