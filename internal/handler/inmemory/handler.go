package inmemory

import (
	"encoding/json"
	"net/http"

	"github.com/Quaqmre/go-case/internal/handler"
	"github.com/Quaqmre/go-case/internal/model"
	"github.com/Quaqmre/go-case/internal/model/inmemory"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
)

type MemStore interface {
	Set(key, value string) error
	Get(key string) (string, error)
}

type Handler struct {
	db MemStore
}

func NewHandler(database MemStore) *Handler {
	return &Handler{db: database}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rb := handler.NewResponseBuilder(w)

	switch r.Method {
	case http.MethodPost:
		h.post(w, r)

	case http.MethodGet:
		h.get(w, r)

	default:
		err := errors.New("not found")
		rb.JsonResponse(http.StatusNotFound, model.NewErrorResponse(err))
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	var keyValue inmemory.Request
	rb := handler.NewResponseBuilder(w)

	err := json.NewDecoder(r.Body).Decode(&keyValue)
	if err != nil {
		rb.JsonResponse(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	validate := validator.New()
	err = validate.Struct(keyValue)
	if err != nil {
		rb.JsonResponse(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	err = h.db.Set(keyValue.Key, keyValue.Value)
	if err != nil {
		rb.JsonResponse(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	rb.JsonResponse(http.StatusOK, keyValue)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	rb := handler.NewResponseBuilder(w)

	if key == "" {
		err := errors.New("key not given")
		rb.JsonResponse(http.StatusBadRequest, model.NewErrorResponse(err))
	}

	value, err := h.db.Get(key)
	if err != nil {
		rb.JsonResponse(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	response := &inmemory.Response{
		Key:   key,
		Value: value,
	}

	rb.JsonResponse(http.StatusOK, response)

}
