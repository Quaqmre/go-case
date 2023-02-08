package handler

import (
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ResponseBuilder struct {
	responseWriter http.ResponseWriter
}

func NewResponseBuilder(responseWriter http.ResponseWriter) *ResponseBuilder {
	return &ResponseBuilder{responseWriter: responseWriter}
}

func (rb *ResponseBuilder) jsonResponseRaw(statusCode int, data []byte) {
	rb.responseWriter.Header().Set("Content-Type", "application/json")
	rb.responseWriter.WriteHeader(statusCode)

	_, err := rb.responseWriter.Write(data)
	if err != nil {
		log.Errorln(errors.Wrapf(err, "could not write data to reponse: %v", data))
	}
}

func (rb *ResponseBuilder) JsonResponse(statusCode int, value interface{}) {
	data, err := json.Marshal(value)
	if err != nil {
		log.Errorln(errors.Wrapf(err, "could not marshal data: %v", value))
	}

	rb.jsonResponseRaw(statusCode, data)
}
