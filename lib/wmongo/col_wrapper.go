package wmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionWrapper interface {
	GenSortBson(sort []string) (result bson.D)
	FindOneAndUpdate(ctx context.Context, filter, update, result interface{}, sort []string, upsert, returnNew bool, opts ...*options.FindOneAndUpdateOptions) (has bool, err error)
	Find(ctx context.Context, filter interface{}, result interface{}, sort []string, skip, limit int64, opts ...*options.FindOptions) (err error)
	FindOne(ctx context.Context, filter interface{}, result interface{}, sort []string, skip int64, opts ...*options.FindOneOptions) (has bool, err error)
	FindById(ctx context.Context, Id interface{}, result interface{}, opts ...*options.FindOneOptions) (has bool, err error)

	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (insertedID interface{}, err error)

	UpdateOne(ctx context.Context, filter, update interface{}, upsert bool, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error)
	Upsert(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error)
	UpsertID(ctx context.Context, ID interface{}, update interface{}) (result *mongo.UpdateResult, err error)

	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (has bool, err error)
	DeleteID(ctx context.Context, ID interface{}, opts ...*options.DeleteOptions) (has bool, err error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (deletedCnt int64, err error)
}
