package usecase

import (
	"database/sql"
	"forex/exchange"

	"forex/models"
)

type exchangeUsecase struct {
	exchangeRepo exchange.ExchangeRepository
}

// create instance
func CreateExchangeUsecase(repo exchange.ExchangeRepository) exchange.ExchangeUsecase {
	return &exchangeUsecase{
		exchangeRepo: repo,
	}
}

// use case to get all data exchange
func (u *exchangeUsecase) Fetch(page int, size int) ([]*models.Exchange, error) {
	if page == 1 {
		page = 0
	}
	if size == 0 {
		size = 10
	}
	return u.exchangeRepo.Fetch(page, size)
}

// use case to create data
func (u *exchangeUsecase) Create(data *models.Exchange) (*models.Exchange, error) {
	_, err := u.exchangeRepo.GetByFromTo(data.From, data.To, false)
	if err != nil && err == sql.ErrNoRows {
		// if the data doesn't exist, insert the data
		lastInsertedId, err := u.exchangeRepo.Create(data)
		if err != nil {
			return nil, err
		}
		data.Id = lastInsertedId
		return data, nil
	}
	return nil, models.CONFLIT_ERROR
}

// use case to delete data by from and to attribute
func (u *exchangeUsecase) Destroy(from string, to string) error {
	existData, err := u.exchangeRepo.GetByFromTo(from, to, true)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.NOT_FOUND_ERROR
		}
		return err
	}
	err = u.exchangeRepo.Destroy(existData)
	if err != nil {
		return err
	}
	return nil
}
