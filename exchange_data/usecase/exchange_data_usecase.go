package usecase

import (
	"database/sql"
	"app/exchange"
	"app/exchange_data"
	"app/helper"

	"app/models"
)

type exchangeDataUseCase struct {
	exchangeDataRepo exchange_data.ExchangeDataRepository
	_exchangeRepo    exchange.ExchangeRepository
}

func CreateExchangeDataUseCase(dataRepo exchange_data.ExchangeDataRepository, exRepo exchange.ExchangeRepository) exchange_data.ExchangeDataUseCase {
	return &exchangeDataUseCase{
		exchangeDataRepo: dataRepo,
		_exchangeRepo:    exRepo,
	}
}

func (u *exchangeDataUseCase) Create(data *models.ExchangeDataRequest) (*models.ExchangeDataRequest, error) {
	// parse Time
	parsedTime, err := helper.ParseTime(data.Date)
	if err != nil {
		return nil, models.INVALID_REQUEST_ERROR
	}

	// check exchange
	currency, err := u._exchangeRepo.GetByFromTo(data.From, data.To, true)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, models.NOT_FOUND_ERROR
		default:
			return nil, err
		}
	}

	// insert data to the database
	exchangeData := &models.ExchangeData{}
	exchangeData.Date = parsedTime
	exchangeData.ExchangeId = currency.Id
	exchangeData.Rate = data.Rate
	resultId, err := u.exchangeDataRepo.Create(exchangeData)

	data.Id = resultId
	return data, nil
}

func (u *exchangeDataUseCase) GetAvgRate(param string) ([]*models.ExchangeAvgAggregate, error) {
	// get time
	parsedTime, err := helper.ParseTime(param)
	if err != nil {
		return nil, models.INVALID_REQUEST_ERROR
	}
	// set start time to 7 days before requested date
	startTime := parsedTime.AddDate(0, -7, 0)

	result, err := u.exchangeDataRepo.GetAvgDay(startTime, parsedTime)
	if err != nil {
		return nil, err
	}
	return result, nil
}
