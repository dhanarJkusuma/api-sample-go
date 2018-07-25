package exchange

import (
	"forex/models"
)

type ExchangeRepository interface {
	Fetch(page int, size int) ([]*models.Exchange, error)
	Create(data *models.Exchange) (int64, error)
	GetByFromTo(from string, to string, handleErrorLog bool) (*models.Exchange, error)
	Destroy(data *models.Exchange) error
}
