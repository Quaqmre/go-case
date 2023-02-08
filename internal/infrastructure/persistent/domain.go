package persistent

import "time"

type DataQuery struct {
	StartDate time.Time
	EndDate   time.Time
	MinCount  int
	MaxCount  int
}

type DataQueryRecord struct {
	Key        string    `bson:"key"`
	CreatedAt  time.Time `bson:"createdAt"`
	TotalCount int       `bson:"totalCount"`
}
