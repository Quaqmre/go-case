package persistent

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Query interface {
	Get(dq *DataQuery) ([]DataQueryRecord, error)
	Disconnect() error
}

type MongoClient interface {
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Disconnect(ctx context.Context) error
}
type MongoCollection interface {
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
}

type Db struct {
	ctx        context.Context
	client     MongoClient
	collection MongoCollection
}

func NewDb(ctx context.Context, mongoClient MongoClient, mongoCollection MongoCollection) (Query, error) {
	db := &Db{
		ctx:        ctx,
		client:     mongoClient,
		collection: mongoCollection,
	}
	return db, nil
}

func (d *Db) Get(dq *DataQuery) ([]DataQueryRecord, error) {
	client := d.client
	collection := d.collection

	if err := client.Ping(d.ctx, readpref.Primary()); err != nil {
		return nil, errors.Wrap(err, "unable to ping the primary")
	}

	agregatedData, err := collection.Aggregate(d.ctx, QueryData(dq))
	if err != nil {
		return nil, errors.Wrap(err, "unexpected behavior in query")
	}

	var withInfo []DataQueryRecord

	if err = agregatedData.All(d.ctx, &withInfo); err != nil {
		return nil, errors.Wrap(err, "unable to iterate over the results")
	}

	return withInfo, nil
}

func (d *Db) Disconnect() error {
	err := d.client.Disconnect(d.ctx)
	if err != nil {
		return errors.Wrap(err, "database disconnecting error")
	}

	return nil
}
