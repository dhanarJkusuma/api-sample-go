package repository

import (
	"database/sql"

	_ "github.com/lib/pq"

	"app/exchange"
	"app/models"
	"log"
)

// Inner Repository Struct
type psqlExchangeRepository struct {
	Conn *sql.DB
}

func CreatePsqlExchangeRespository(Conn *sql.DB) exchange.ExchangeRepository {
	return &psqlExchangeRepository{Conn}
}

// private function fetch
func (p *psqlExchangeRepository) fetch(query string, params ...interface{}) ([]*models.Exchange, error) {
	rows, err := p.Conn.Query(query, params...)
	if err != nil {
		log.Printf("[Error DB] : %v", err)
		return nil, err
	}

	exchangeList := make([]*models.Exchange, 0)
	for rows.Next() {
		exchange := &models.Exchange{}
		err := rows.Scan(&exchange.Id, &exchange.From, &exchange.To)
		if err != nil {
			log.Printf("[Error DB] : %v", err)
			return nil, err
		}
		exchangeList = append(exchangeList, exchange)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[Error DB] : %v", err)
		return nil, err
	}
	return exchangeList, nil
}

// private function to get single object
func (p *psqlExchangeRepository) get(handleErrorLog bool, query string, params ...interface{}) (*models.Exchange, error) {
	row := p.Conn.QueryRow(query, params...)
	exchange := &models.Exchange{}
	err := row.Scan(&exchange.Id, &exchange.From, &exchange.To)
	if err != nil {
		if handleErrorLog {
			log.Printf("[Error DB] : %v", err)
		}
		return nil, err
	}
	return exchange, nil
}

// public function to get data by from & to attribute
func (p *psqlExchangeRepository) GetByFromTo(from string, to string, handleErrorLog bool) (*models.Exchange, error) {
	query := "SELECT * FROM exchange_rate WHERE from_cur=$1 AND to_cur=$2"
	exchange, err := p.get(handleErrorLog, query, from, to)
	if err != nil {
		return nil, err
	}
	return exchange, nil
}

// public function to create exchange
func (p *psqlExchangeRepository) Create(data *models.Exchange) (int64, error) {
	var lastId int64
	query := "INSERT INTO exchange_rate (from_cur, to_cur) VALUES ($1, $2) RETURNING id"
	_ = p.Conn.QueryRow(query, data.From, data.To).Scan(&lastId)
	return lastId, nil
}

// public function to fetch exchange data
func (p *psqlExchangeRepository) Fetch(page int, size int) ([]*models.Exchange, error) {
	offset := page * size
	query := "SELECT * FROM exchange_rate LIMIT $1 OFFSET $2;"
	return p.fetch(query, size, offset)
}

func (p *psqlExchangeRepository) Destroy(data *models.Exchange) error {
	query := "DELETE FROM exchange_rate WHERE id=$1"
	_, err := p.Conn.Exec(query, data.Id)
	if err != nil {
		return err
	}
	return nil
}
