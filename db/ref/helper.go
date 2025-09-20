package ref

import "go.mongodb.org/mongo-driver/v2/bson"

type IMongoHelper interface {
}

type MongoHelper struct {
}

func NewMongoHelper() IMongoHelper {
	return &MongoHelper{}
}

func UpdateSet(update any) any {
	return bson.M{"$set": update}
}

func UpdateUnset(update any) any {
	return bson.M{"$unset": update}
}

func UpdateSetPipeline(update any) any {
	return []bson.M{{"$set": update}}
}

// FindOption allows customizing find operations
type FindOption func(*FindOptions)

type FindOptions struct {
	Limit      *int64
	Skip       *int64
	Sort       any
	Projection any
}

// WithLimit sets the limit for find operations
func WithLimit(limit int64) FindOption {
	return func(opts *FindOptions) {
		opts.Limit = &limit
	}
}

// WithSort sets the sort order for find operations
func WithSort(sort any) FindOption {
	return func(opts *FindOptions) {
		opts.Sort = sort
	}
}

// WithSkip sets the number of documents to skip
func WithSkip(skip int64) FindOption {
	return func(opts *FindOptions) {
		opts.Skip = &skip
	}
}

// WithProjection sets which fields to include/exclude in the result
func WithProjection(projection any) FindOption {
	return func(opts *FindOptions) {
		opts.Projection = projection
	}
}
