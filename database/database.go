package database

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Database interface {
	Collection(string) Collection
	Client() Client
}

type Collection interface {
	FindOne(context.Context, interface{}) SingleResult
	Find(context.Context, interface{}, ...options.Lister[options.FindOptions]) (Cursor, error)
	InsertOne(context.Context, interface{}) (string, error)
	DeleteOne(context.Context, interface{}) (int64, error)
	UpdateOne(context.Context, interface{}, interface{}, ...options.Lister[options.UpdateOptions]) (UpdateResult, error)
	UpdateMany(context.Context, interface{}, interface{}, ...options.Lister[options.UpdateOptions]) (UpdateResult, error)
	ReplaceOne(context.Context, interface{}, interface{}, ...options.Lister[options.ReplaceOptions]) (UpdateResult, error)
}

type SingleResult interface {
	Decode(interface{}) error
}

type UpdateResult struct {
	MatchedCount  int64
	ModifiedCount int64
}

type FindOptions struct {
	Limit  int64
	Skip   int64
	SortBy string
}

type Cursor interface {
	Close(context.Context) error
	Next(context.Context) bool
	Decode(interface{}) error
	All(context.Context, interface{}) error
}

type Client interface {
	Database(string) Database
	Ping(context.Context) error
	Disconnect(context.Context) error
}

type mongoClient struct {
	cl *mongo.Client
}
type mongoDatabase struct {
	db *mongo.Database
}
type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoUpdateResult struct {
	ur *mongo.UpdateResult
}

func (mc *mongoClient) Ping(ctx context.Context) error {
	return mc.cl.Ping(ctx, readpref.Primary())
}

func (mc *mongoClient) Database(dbName string) Database {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (md *mongoDatabase) Collection(colName string) Collection {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() Client {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResult {
	singleResult := mc.coll.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOptions]) (Cursor, error) {
	cursor, err := mc.coll.Find(ctx, filter, opts[:]...)
	return cursor, err
}

func (mc *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...options.Lister[options.UpdateOptions]) (UpdateResult, error) {
	res, err := mc.coll.UpdateOne(ctx, filter, update, opts[:]...)
	if err != nil {
		return UpdateResult{}, err
	}

	ur := UpdateResult{
		MatchedCount:  res.MatchedCount,
		ModifiedCount: res.ModifiedCount,
	}
	return ur, nil
}

func (mc *mongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...options.Lister[options.UpdateOptions]) (UpdateResult, error) {
	res, err := mc.coll.UpdateMany(ctx, filter, update, opts[:]...)
	if err != nil {
		return UpdateResult{}, err
	}

	ur := UpdateResult{
		MatchedCount:  res.MatchedCount,
		ModifiedCount: res.ModifiedCount,
	}
	return ur, nil
}

func (mc *mongoCollection) ReplaceOne(ctx context.Context, filter interface{}, update interface{}, opts ...options.Lister[options.ReplaceOptions]) (UpdateResult, error) {
	res, err := mc.coll.ReplaceOne(ctx, filter, update, opts[:]...)
	if err != nil {
		return UpdateResult{}, err
	}

	ur := UpdateResult{
		MatchedCount:  res.MatchedCount,
		ModifiedCount: res.ModifiedCount,
	}
	return ur, nil
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (string, error) {
	id, err := mc.coll.InsertOne(ctx, document)

	if id == nil {
		return "", err
	}
	return id.InsertedID.(string), err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.coll.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

func (mc *mongoClient) Disconnect(ctx context.Context) error {
	return mc.cl.Disconnect(ctx)
}

func NewClient(connection string) (Client, error) {
	c, err := mongo.Connect(options.Client().ApplyURI(connection))
	return &mongoClient{cl: c}, err
}
