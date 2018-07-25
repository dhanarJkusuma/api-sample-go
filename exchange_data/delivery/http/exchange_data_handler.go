package http

import (
	goHttp "net/http"

	"app/helper"

	"app/models"

	"app/exchange_data"

	"github.com/julienschmidt/httprouter"
	validator "gopkg.in/go-playground/validator.v9"
)

type HttpExchangeDataHandler struct {
	UseCase exchange_data.ExchangeDataUseCase
}

func (h *HttpExchangeDataHandler) Create(w goHttp.ResponseWriter, req *goHttp.Request, _ httprouter.Params) {
	helper.HandleCORS(&w, req)

	var data models.ExchangeDataRequest
	helper.UnMarshall(req.Body, &data)

	if ok, err := isRequestValid(&data); !ok {
		helper.RespondError(w, goHttp.StatusBadRequest, err.Error())
		return
	}

	result, err := h.UseCase.Create(&data)
	if err != nil {
		switch err {
		case models.NOT_FOUND_ERROR:
			helper.RespondError(w, goHttp.StatusNotFound, "Exchange currency not found.")
			return
		case models.INVALID_REQUEST_ERROR:
			helper.RespondError(w, goHttp.StatusBadRequest, "Invalid date format. Format (yyyy-MM-dd) Required.")
			return
		default:
			helper.RespondError(w, goHttp.StatusInternalServerError, "Internal Server Error.")
			return
		}
	}

	helper.RespondJSON(w, goHttp.StatusCreated, result)
}

func (h *HttpExchangeDataHandler) GetAvg(w goHttp.ResponseWriter, req *goHttp.Request, _ httprouter.Params) {
	helper.HandleCORS(&w, req)

	query := req.URL.Query()
	date := query.Get("search_date")
	if date == "" {
		helper.RespondError(w, goHttp.StatusBadRequest, `Query "search_date" is Required.`)
		return
	}
	result, err := h.UseCase.GetAvgRate(date)
	if err != nil {
		switch err {
		case models.INVALID_REQUEST_ERROR:
			helper.RespondError(w, goHttp.StatusBadRequest, "Invalid date format. Format (yyyy-MM-dd) Required.")
			return
		default:
			helper.RespondError(w, goHttp.StatusInternalServerError, "Internal Server Error.")
			return
		}
	}
	helper.RespondJSON(w, goHttp.StatusOK, result)
}

func isRequestValid(m *models.ExchangeDataRequest) (bool, error) {
	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CreateHttpExchangeDataHandler(router *httprouter.Router, uc exchange_data.ExchangeDataUseCase) {
	handler := &HttpExchangeDataHandler{
		UseCase: uc,
	}

	router.POST("/api/exchange/rate", handler.Create)
	router.GET("/api/exchange/rate", handler.GetAvg)
}
