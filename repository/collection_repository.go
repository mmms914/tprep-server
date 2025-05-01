package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"main/database"
	"main/domain"
	"main/internal"
)

type collectionRepository struct {
	database   database.Database
	collection string
}

func NewCollectionRepository(db database.Database, collection string) domain.CollectionRepository {
	return &collectionRepository{
		database:   db,
		collection: collection,
	}
}

func (cr *collectionRepository) Create(c context.Context, collection *domain.Collection) (string, error) {
	collections := cr.database.Collection(cr.collection)
	collection.ID = internal.GenerateUUID()
	id, err := collections.InsertOne(c, collection)
	return id, err
}

func (cr *collectionRepository) UpdateByID(c context.Context, collectionID string, update interface{}) (database.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: collectionID}}
	return cr.Update(c, filter, update)
}

func (cr *collectionRepository) Update(c context.Context, filter interface{}, update interface{}) (database.UpdateResult, error) {
	collections := cr.database.Collection(cr.collection)
	return collections.UpdateOne(c, filter, update)
}

func (cr *collectionRepository) DeleteByID(c context.Context, collectionID string) error {
	collections := cr.database.Collection(cr.collection)
	filter := bson.D{{Key: "_id", Value: collectionID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "is_deleted", Value: true}}}}
	_, err := collections.UpdateOne(c, filter, update)
	return err
}

func (cr *collectionRepository) GetByID(c context.Context, collectionID string) (domain.Collection, error) {
	var result domain.Collection
	collections := cr.database.Collection(cr.collection)
	filter := bson.D{
		{Key: "_id", Value: collectionID},
		{Key: "is_deleted", Value: bson.D{
			{Key: "$exists", Value: false},
		}},
	}
	err := collections.FindOne(c, filter).Decode(&result)
	return result, err
}

func (cr *collectionRepository) GetByFilter(c context.Context, filter interface{}, opts database.FindOptions) ([]domain.Collection, error) {
	var results []domain.Collection
	collections := cr.database.Collection(cr.collection)

	notDeletedFilter := bson.D{
		{"$and", []interface{}{
			filter,
			bson.D{{Key: "is_deleted", Value: bson.D{{Key: "$exists", Value: false}}}},
		}},
	}

	op := options.Find()
	if opts.SortBy != "" {
		sortFilter := bson.D{
			{Key: opts.SortBy, Value: -1},
			{Key: "name", Value: 1},
			{Key: "_id", Value: 1},
		}
		op = op.SetSort(sortFilter)
	}
	op = op.SetLimit(opts.Limit).SetSkip(opts.Skip)

	cursor, err := collections.Find(c, notDeletedFilter, op)
	if err != nil {
		return nil, err
	}
	err = cursor.All(c, &results)
	if err != nil {
		return nil, err
	}
	return results, err
}
