package repository_test

import (
	"testing"
	"time"

	exchangeRateRepo "app/exchange_data/repository"
	"app/models"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// Testing Success Method Create
func TestCreate(t *testing.T) {

	// init mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("[Forex Test] : Error while opening connection. Caused : %s", err)
	}
	defer db.Close()

	// init request record
	ed := &models.ExchangeData{
		ExchangeId: 1,
		Date:       time.Now(),
		Rate:       0.61424,
	}

	// mock result obj
	row := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// mock query
	query := "INSERT INTO exchange_rate_data (.+) RETURNING id"

	// mock result
	mock.ExpectQuery(query).WithArgs(ed.ExchangeId, ed.Date, ed.Rate).WillReturnRows(row)

	// init mock repo
	edr := exchangeRateRepo.CreateExchangeDataRepository(db)

	// insert
	lastId, err := edr.Create(ed)

	// assert Error
	assert.NoError(t, err)

	// assert Result
	assert.Equal(t, int64(1), lastId)

}

// Testing Success Method Avg
func TestGetAvgDay(t *testing.T) {

	// init mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("[Forex Test] : Error while opening connection. Caused : %s", err)
	}
	defer db.Close()

	// mock result obj
	rows := sqlmock.NewRows([]string{"id", "date", "rate", "from_cur", "to_cur", "average"}).
		AddRow(1, time.Now(), 0.6374561, "USD", "JPY", 0.735281).
		AddRow(1, time.Now(), 0.6374561, "USD", "JPY", 0.735281)

	query := `WITH find_record_rate`

	mock.ExpectQuery(query).WillReturnRows(rows)

	// init mock repo
	edr := exchangeRateRepo.CreateExchangeDataRepository(db)

	start := time.Now().AddDate(0, -7, 0)
	result, err := edr.GetAvgDay(start, time.Now())

	// assert Error
	assert.NoError(t, err)

	// assert Length
	assert.Len(t, result, 2)

}
