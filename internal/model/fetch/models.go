package fetch

import "time"

type Request struct {
	StartDate string `json:"startDate" validate:"required"`
	EndDate   string `json:"endDate" validate:"required"`
	MinCount  int    `json:"minCount" validate:"required"`
	MaxCount  int    `json:"maxCount" validate:"required"`
}

type RecordResponse struct {
	Key        string    `json:"key"`
	CreatedAt  time.Time `json:"createdAt"`
	TotalCount int       `json:"totalCount"`
}

type Response struct {
	Code    int              `json:"code"`
	Msg     string           `json:"msg"`
	Records []RecordResponse `json:"records"`
}

func NewErrorResponse(err error) *Response {
	return &Response{
		Code:    1,
		Msg:     err.Error(),
		Records: nil,
	}
}
