package util

import (
	"time"

	"github.com/Quaqmre/go-case/internal/infrastructure/persistent"
	"github.com/Quaqmre/go-case/internal/model/fetch"
	"github.com/pkg/errors"
)

func RequestToDataQuery(model *fetch.Request) (*persistent.DataQuery, error) {
	startDate, err := time.Parse("2006-01-02", model.StartDate)
	if err != nil {
		return nil, errors.Wrap(err, "invalid date format")
	}

	endDate, err := time.Parse("2006-01-02", model.EndDate)
	if err != nil {
		return nil, errors.Wrap(err, "invalid date format")
	}

	if startDate.After(endDate) {
		return nil, errors.New("dates are not in correct chronological order")
	}

	return &persistent.DataQuery{
		StartDate: startDate,
		EndDate:   endDate,
		MinCount:  model.MinCount,
		MaxCount:  model.MaxCount,
	}, nil
}

func RecordsToResponses(records []persistent.DataQueryRecord) []fetch.RecordResponse {
	responses := make([]fetch.RecordResponse, 0)

	for _, record := range records {
		response := fetch.RecordResponse{
			Key:        record.Key,
			CreatedAt:  record.CreatedAt,
			TotalCount: record.TotalCount,
		}

		responses = append(responses, response)
	}

	return responses
}
