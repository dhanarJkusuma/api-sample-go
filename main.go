package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	exchangeHttp "app/exchange/delivery/http"
	exchangeRepo "app/exchange/repository"
	exchangeUsecase "app/exchange/usecase"

	exchangeDataHttp "app/exchange_data/delivery/http"
	exchangeDataRepo "app/exchange_data/repository"
	exchangeDataUseCase "app/exchange_data/usecase"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func init() {
	fmt.Println("[Forex App] : Init Config. ")
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
	fmt.Println("[Forex App] : Connecting to the Database. ")
	stringConfig := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbName)
	dbConn, err := sql.Open("postgres", stringConfig)
	if err != nil {
		log.Fatalln(err)
	}
	defer dbConn.Close()

	if err = dbConn.Ping(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("[Forex App] : Database Connected. ")

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

	handlerCors := cors.Default().Handler(router)

	fmt.Println("[Forex App] : Application is running. ")
	http.ListenAndServe(":8080", handlerCors)
}
