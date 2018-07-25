package repository

import (
	"database/sql"
	"log"
	"time"

	"forex/exchange_data"
	"forex/helper"
	"forex/models"
	"strconv"

	_ "github.com/lib/pq"
)

type psqlExchangeDataRepository struct {
	Conn *sql.DB
}

func CreateExchangeDataRepository(Conn *sql.DB) exchange_data.ExchangeDataRepository {
	return &psqlExchangeDataRepository{
		Conn,
	}
}

func (p *psqlExchangeDataRepository) Create(data *models.ExchangeData) (int64, error) {
	var lastId int64
	query := "INSERT INTO exchange_rate_data (exchange_rate_id, date, rate) VALUES ($1, $2, $3)  RETURNING id"
	_ = p.Conn.QueryRow(query, data.ExchangeId, data.Date, data.Rate).Scan(&lastId)
	return lastId, nil
}

func (p *psqlExchangeDataRepository) GetAvgDay(start time.Time, end time.Time) ([]*models.ExchangeAvgAggregate, error) {
	query := `
			WITH find_record_rate AS (
			SELECT
				d.id, d.date, d.exchange_rate_id, d.rate
			FROM exchange_rate_data d
			WHERE date=$2
		), find_average AS (
			SELECT e.id, AVG(d.rate) as average FROM
			exchange_rate_data d
			INNER JOIN exchange_rate e
			ON
			d.exchange_rate_id=e.id
			WHERE d.date BETWEEN $1 AND $2
			GROUP BY e.id
		)
		SELECT distinct on (e.id) e.id, d.date, coalesce(d.rate, -1.00), e.from_cur, e.to_cur, coalesce(av.average, 0.00)
		FROM
		find_record_rate d
		RIGHT JOIN exchange_rate e
		ON
			d.exchange_rate_id=e.id
		LEFT JOIN find_average av
		ON
			av.id=e.id
		GROUP BY e.id, d.date, d.rate, av.average;
	`
	rows, err := p.Conn.Query(query, start, end)
	if err != nil {
		log.Printf("[Error DB] : %v", err)
		return nil, err
	}

	recordList := make([]*models.ExchangeAvgAggregate, 0)
	for rows.Next() {
		ag := &models.ExchangeAvgAggregate{}
		err := rows.Scan(&ag.ExchangeId, &ag.Date, &ag.Rate, &ag.From, &ag.To, &ag.Average)
		if err != nil {
			log.Printf("[Error DB] : %v", err)
			return nil, err
		}
		// parse date
		if ag.Date.Valid {
			ag.DateString = helper.ParseTimeString(ag.Date.Time)
		} else {
			ag.DateString = helper.ParseTimeString(end)
		}

		// parse rate
		if ag.Rate == -1.0 {
			ag.RateString = "insufficient data"
		} else {
			ag.RateString = strconv.FormatFloat(ag.Rate, 'f', 15, 64)
		}
		recordList = append(recordList, ag)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[Error DB] : %v", err)
		return nil, err
	}
	return recordList, nil
}
