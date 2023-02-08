package persistent

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func QueryData(dq *DataQuery) []primitive.M {
	datas := []bson.M{
		{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gt": dq.StartDate,
					"$lt": dq.EndDate,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":        0,
				"key":        1,
				"createdAt":  1,
				"totalCount": bson.M{"$sum": "$counts"},
			},
		},
		{
			"$match": bson.M{
				"totalCount": bson.M{
					"$gt": dq.MinCount,
					"$lt": dq.MaxCount,
				},
			},
		},
	}

	return datas
}
