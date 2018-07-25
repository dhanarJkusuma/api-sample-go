package exchange

import (
	"app/models"
)

type ExchangeUsecase interface {
	Fetch(page int, size int) ([]*models.Exchange, error)
	Create(data *models.Exchange) (*models.Exchange, error)
	Destroy(from string, to string) error
}
