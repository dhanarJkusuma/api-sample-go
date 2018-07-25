package http

import (
	"app/exchange"
	goHttp "net/http"

	"app/helper"
	"strconv"

	"app/models"

	"github.com/julienschmidt/httprouter"
	validator "gopkg.in/go-playground/validator.v9"
)

type HttpExchangeHandler struct {
	UseCase exchange.ExchangeUsecase
}

// handler for fetch
func (h *HttpExchangeHandler) Fetch(w goHttp.ResponseWriter, req *goHttp.Request, params httprouter.Params) {
	helper.HandleCORS(&w, req)
	page, err := strconv.Atoi(params.ByName("page"))
	if err != nil {
		page = 0
	}
	size, err := strconv.Atoi(params.ByName("size"))
	if err != nil {
		size = 0
	}
	result, err := h.UseCase.Fetch(page, size)
	if err != nil {
		helper.RespondError(w, goHttp.StatusInternalServerError, "Internal Server Error")
		return
	}
	helper.RespondJSON(w, goHttp.StatusOK, result)
}

// handler for create
func (h *HttpExchangeHandler) Create(w goHttp.ResponseWriter, req *goHttp.Request, _ httprouter.Params) {
	helper.HandleCORS(&w, req)
	var data models.Exchange
	helper.UnMarshall(req.Body, &data)

	if ok, err := isRequestValid(&data); !ok {
		helper.RespondError(w, goHttp.StatusBadRequest, err.Error())
		return
	}

	result, err := h.UseCase.Create(&data)
	if err != nil {
		switch err {
		case models.CONFLIT_ERROR:
			helper.RespondError(w, goHttp.StatusConflict, err.Error())
			return
		default:
			helper.RespondError(w, goHttp.StatusInternalServerError, "Something were wrong !")
			return
		}
	}

	helper.RespondJSON(w, goHttp.StatusCreated, result)
}

// handler for delete
func (h *HttpExchangeHandler) Destroy(w goHttp.ResponseWriter, req *goHttp.Request, params httprouter.Params) {
	helper.HandleCORS(&w, req)
	query := req.URL.Query()
	from := query.Get("from")
	if from == "" {
		helper.RespondError(w, goHttp.StatusBadRequest, `Query "from" is Required.`)
		return
	}

	to := query.Get("to")
	if to == "" {
		helper.RespondError(w, goHttp.StatusBadRequest, `Query "to" is Required.`)
		return
	}

	err := h.UseCase.Destroy(from, to)
	if err != nil {
		switch err {
		case models.NOT_FOUND_ERROR:
			helper.RespondError(w, goHttp.StatusNotFound, err.Error())
			return
		default:
			helper.RespondError(w, goHttp.StatusInternalServerError, "Something were wrong !")
			return
		}
	}
	helper.RespondJSON(w, goHttp.StatusNoContent, nil)
}

func isRequestValid(m *models.Exchange) (bool, error) {
	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CreateHttpExchangeHandler(router *httprouter.Router, uc exchange.ExchangeUsecase) {
	handler := &HttpExchangeHandler{
		UseCase: uc,
	}

	router.GET("/api/exchange", handler.Fetch)
	router.POST("/api/exchange", handler.Create)
	router.DELETE("/api/exchange", handler.Destroy)

}
