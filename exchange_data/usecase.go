package exchange_data

import "forex/models"

type ExchangeDataUseCase interface {
	Create(data *models.ExchangeDataRequest) (*models.ExchangeDataRequest, error)
	GetAvgRate(param string) ([]*models.ExchangeAvgAggregate, error)
}
