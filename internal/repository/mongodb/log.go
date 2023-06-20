package mongodb

import (
	"context"

	errwrap "github.com/pkg/errors"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/helper"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/repository/mongodb/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogRepository interface {
	Create(ctx context.Context, params entity.LogCollection) error
}

type Log struct {
	collection *mongo.Collection
}

func NewLogRepository(db *mongo.Database) *Log {
	return &Log{collection: db.Collection(LogCollection)}
}

func (r *Log) Create(ctx context.Context, params entity.LogCollection) error {
	funcName := "[LogRepositoryMongo.Create]"

	if err := helper.CheckDeadline(ctx); err != nil {
		return errwrap.Wrap(err, funcName)
	}

	_, err := r.collection.InsertOne(ctx, params)
	return err
}
