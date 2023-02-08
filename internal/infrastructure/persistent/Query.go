package persistent

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func QueryData(dq *DataQuery) mongo.Pipeline {
	filterCreateDate := bson.D{{"$match", bson.D{
		{"createdAt", bson.D{
			{"$gte", dq.StartDate},
			{"$lt", dq.EndDate},
		}}}}}

	project := bson.D{{"$project", bson.D{
		{"_id", false},
		{"key", true},
		{"createdAt", true},
		{"totalCount", bson.D{
			{"$sum", "$counts"},
		}}}}}

	filterTotalCount := bson.D{{"$match", bson.D{
		{"totalCount", bson.D{
			{"$gte", dq.MinCount},
			{"$lt", dq.MaxCount},
		}},
	}}}

	return mongo.Pipeline{filterCreateDate, project, filterTotalCount}
}
