package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	exchangeHttp "forex/exchange/delivery/http"
	exchangeRepo "forex/exchange/repository"
	exchangeUsecase "forex/exchange/usecase"

	exchangeDataHttp "forex/exchange_data/delivery/http"
	exchangeDataRepo "forex/exchange_data/repository"
	exchangeDataUseCase "forex/exchange_data/usecase"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}
}

func main() {

	// get config from config.json
	host := viper.GetString(`database.host`)
	port := viper.GetString(`database.port`)
	user := viper.GetString(`database.user`)
	pass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	// open db connection
	stringConfig := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbName)
	dbConn, err := sql.Open("postgres", stringConfig)
	if err != nil {
		log.Fatalln(err)
	}
	defer dbConn.Close()

	if err = dbConn.Ping(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Database connected successfully.")

	// init router
	router := httprouter.New()

	// init repository
	er := exchangeRepo.CreatePsqlExchangeRespository(dbConn)
	edr := exchangeDataRepo.CreateExchangeDataRepository(dbConn)

	// init usecase
	eu := exchangeUsecase.CreateExchangeUsecase(er)
	edu := exchangeDataUseCase.CreateExchangeDataUseCase(edr, er)

	// init delivery/http
	exchangeHttp.CreateHttpExchangeHandler(router, eu)
	exchangeDataHttp.CreateHttpExchangeDataHandler(router, edu)

	http.ListenAndServe(":8080", router)

}
