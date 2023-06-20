package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongodb(ctx context.Context, cfg *MongodbOption) (*mongo.Database, error) {
	opts := mongoOptions(cfg)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// check connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.DatabaseName), nil
}

func mongoOptions(cfg *MongodbOption) *options.ClientOptions {
	return options.Client().ApplyURI(cfg.Uri)
}
