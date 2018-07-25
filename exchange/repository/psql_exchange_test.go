package repository_test

import (
	"testing"

	exchangeRepo "app/exchange/repository"
	"app/models"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// Testing Success Method Fetch
func TestFetch(t *testing.T) {
	// init mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("[Forex Test] : Error while opening connection. Caused : %s", err)
	}
	defer db.Close()

	// mocking record
	rows := sqlmock.NewRows([]string{"id", "from_cur", "to_cur"}).
		AddRow(1, "USD", "JPY").
		AddRow(2, "IDR", "JPY")

	query := `SELECT (.+) FROM exchange_rate`

	// mocking query
	mock.ExpectQuery(query).WillReturnRows(rows)
	// init mock repo
	er := exchangeRepo.CreatePsqlExchangeRespository(db)

	// get data
	list, err := er.Fetch(0, 10)

	// assert Error
	assert.NoError(t, err)
	// assert Length
	assert.Len(t, list, 2)
}

// Testing Success Method GetByFromTo
func TestGetByFromTo(t *testing.T) {
	// init mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("[Forex Test] : Error while opening connection. Caused : %s", err)
	}
	defer db.Close()

	// mocking record result
	rows := sqlmock.NewRows([]string{"id", "from_cur", "to_cur"}).
		AddRow(1, "USD", "JPY")

	query := `SELECT (.+) FROM exchange_rate WHERE`

	// mocking query
	mock.ExpectQuery(query).WillReturnRows(rows)
	// init mock repo
	er := exchangeRepo.CreatePsqlExchangeRespository(db)

	// get data
	exchangeCur, err := er.GetByFromTo("USD", "JPY", true)
	// assert Error
	assert.NoError(t, err)
	// assert Length
	assert.NotNil(t, exchangeCur)
}

//  Testing success Method Create
func TestCreate(t *testing.T) {
	// init request record
	ec := &models.Exchange{
		From: "USD",
		To:   "JPY",
	}

	// init mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("[Forex Test] : Error while opening connection. Caused : %s", err)
	}

	defer db.Close()

	// mock result obj
	row := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// mock query
	query := "INSERT INTO exchange_rate (.+) RETURNING id"

	// mock result
	mock.ExpectQuery(query).WithArgs(ec.From, ec.To).WillReturnRows(row)

	// init mock repositorlly
	er := exchangeRepo.CreatePsqlExchangeRespository(db)

	// insert record
	lastId, err := er.Create(ec)

	// assert Error
	assert.NoError(t, err)
	// assert Result
	assert.Equal(t, int64(1), lastId)
}

// Testing success Method Delete
func TestDestroy(t *testing.T) {
	// init mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("[Forex Test] : Error while opening connection. Caused : %s", err)
	}
	defer db.Close()

	// init delete record
	data := &models.Exchange{
		Id:   12,
		From: "USD",
		To:   "JPY",
	}

	// mock query
	query := `DELETE FROM exchange_rate`

	// mock result
	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(12, 1))

	// init mock repositorlly
	er := exchangeRepo.CreatePsqlExchangeRespository(db)

	// delete record
	err = er.Destroy(data)

	// asssert error
	assert.NoError(t, err)

}
