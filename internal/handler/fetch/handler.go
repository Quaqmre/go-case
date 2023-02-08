package fetch

import (
	"encoding/json"
	"net/http"

	"github.com/Quaqmre/go-case/internal/handler"
	"github.com/Quaqmre/go-case/internal/infrastructure/persistent"
	"github.com/Quaqmre/go-case/internal/model"
	"github.com/Quaqmre/go-case/internal/model/fetch"
	"github.com/Quaqmre/go-case/internal/util"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
)

type Database interface {
	Disconnect() error
	Get(dq *persistent.DataQuery) ([]persistent.DataQueryRecord, error)
}

type Handler struct {
	db Database
}

func NewHandler(database Database) *Handler {
	return &Handler{db: database}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rb := handler.NewResponseBuilder(w)

	switch r.Method {
	case http.MethodPost:
		h.post(w, r)
	default:
		err := errors.New("not found")
		rb.JsonResponse(http.StatusNotFound, model.NewErrorResponse(err))
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	request := new(fetch.Request)
	rb := handler.NewResponseBuilder(w)

	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		rb.JsonResponse(http.StatusBadRequest, fetch.NewErrorResponse(err))
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		rb.JsonResponse(http.StatusBadRequest, fetch.NewErrorResponse(err))
		return
	}

	query, err := util.RequestToDataQuery(request)
	if err != nil {
		rb.JsonResponse(http.StatusBadRequest, fetch.NewErrorResponse(err))
		return
	}

	queryRecords, err := h.db.Get(query)
	if err != nil {
		rb.JsonResponse(http.StatusBadRequest, fetch.NewErrorResponse(err))
		return
	}

	responseRecords := util.RecordsToResponses(queryRecords)
	response := &fetch.Response{
		Code:    0,
		Msg:     "Success",
		Records: responseRecords,
	}

	rb.JsonResponse(http.StatusOK, response)
}
