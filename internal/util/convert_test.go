package util_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/Quaqmre/go-case/internal/infrastructure/persistent"
	"github.com/Quaqmre/go-case/internal/model/fetch"
	"github.com/Quaqmre/go-case/internal/util"
)

func TestRequestToDataQuery_With_an_invalid_date_format(t *testing.T) {
	tests := []struct {
		name      string
		input     *fetch.Request
		wantError bool
	}{
		{
			name: "invalid StartDate format",
			input: &fetch.Request{
				StartDate: "23344wer",
				EndDate:   "2018-02-02",
				MinCount:  2000,
				MaxCount:  3000,
			},
			wantError: true,
		},
		{
			name: "invalid EndDate format",
			input: &fetch.Request{
				StartDate: "2016-01-26",
				EndDate:   "20120305",
				MinCount:  2000,
				MaxCount:  3000,
			},
			wantError: true,
		},
		{
			name: "invalid start date",
			input: &fetch.Request{
				StartDate: "2020-34-54",
				EndDate:   "2021-02-02",
				MinCount:  0,
				MaxCount:  0,
			},
			wantError: true,
		},
		{
			name: "invalid end date",
			input: &fetch.Request{
				StartDate: "2016-01-26",
				EndDate:   "2021-56-34",
				MinCount:  0,
				MaxCount:  0,
			},
			wantError: true,
		},
		{
			name: "invalid chronological date order",
			input: &fetch.Request{
				StartDate: "2021-02-02",
				EndDate:   "2018-02-02",
				MinCount:  0,
				MaxCount:  0,
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := util.RequestToDataQuery(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("RequestToDataQuery() error = %v, wantErr %v", err, tt.wantError)
				return
			}
		})
	}
}

func TestRequestToDataQuery_With_an_valid_date_format(t *testing.T) {
	tests := []struct {
		name      string
		input     *fetch.Request
		expected  *persistent.DataQuery
		wantError bool
	}{
		{
			name: "valid date formats",
			input: &fetch.Request{
				StartDate: "2016-01-26",
				EndDate:   "2018-02-02",
				MinCount:  0,
				MaxCount:  0,
			},
			expected: &persistent.DataQuery{
				StartDate: time.Date(2016, 01, 26, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2018, 02, 02, 0, 0, 0, 0, time.UTC),
				MinCount:  0,
				MaxCount:  0,
			},
		},
		{
			name: "consistent counts",
			input: &fetch.Request{
				StartDate: "2020-03-10",
				EndDate:   "2021-05-04",
				MinCount:  1234,
				MaxCount:  8765,
			},
			expected: &persistent.DataQuery{
				StartDate: time.Date(2020, 3, 10, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
				MinCount:  1234,
				MaxCount:  8765,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.RequestToDataQuery(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("RequestToDataQuery() error = %v, wantErr %v", err, tt.wantError)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("RequestToDataQuery() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRequestToDataQuery_With_an_valid_records(t *testing.T) {
	tests := []struct {
		name     string
		input    []persistent.DataQueryRecord
		expected []fetch.RecordResponse
	}{
		{
			name: "unordered valid data",
			input: []persistent.DataQueryRecord{
				{
					Key:        "a",
					CreatedAt:  time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC),
					TotalCount: 1234,
				},
				{
					Key:        "b",
					CreatedAt:  time.Date(2021, 5, 3, 0, 0, 0, 0, time.UTC),
					TotalCount: 9876,
				},
			},
			expected: []fetch.RecordResponse{
				{
					Key:        "b",
					CreatedAt:  time.Date(2021, 5, 3, 0, 0, 0, 0, time.UTC),
					TotalCount: 9876,
				},
				{
					Key:        "a",
					CreatedAt:  time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC),
					TotalCount: 1234,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.RecordsToResponses(tt.input)

			if got[0] != tt.expected[1] {
				t.Errorf("RequestToDataQuery() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
