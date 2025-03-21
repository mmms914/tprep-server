package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"main/database"
	"main/domain"
)

type userHistoryRepository struct {
	database   database.Database
	collection string
}

func NewUserHistoryRepository(db database.Database, collection string) domain.UserHistoryRepository {
	return &userHistoryRepository{
		database:   db,
		collection: collection,
	}
}

func (uhr *userHistoryRepository) CreateIfNotExists(c context.Context, userID string) error {
	collection := uhr.database.Collection(uhr.collection)

	var userHistory domain.UserHistory
	filter := bson.M{"_id": userID}
	if err := collection.FindOne(c, filter).Decode(&userHistory); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			userHistory := domain.UserHistory{
				UserID: userID,
				Items:  make([]domain.HistoryItem, 0),
			}
			_, err = collection.InsertOne(c, userHistory)
			return err
		}
		return err
	}
	return nil
}

func (uhr *userHistoryRepository) UpdateByID(c context.Context, userID string, item domain.HistoryItem) error {
	err := uhr.CreateIfNotExists(c, userID)
	if err != nil {
		return err
	}

	collection := uhr.database.Collection(uhr.collection)
	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.D{
		{"$push", bson.D{
			{"items", item},
		}},
	}
	_, err = collection.UpdateOne(c, filter, update)

	return err
}

func (uhr *userHistoryRepository) GetByID(c context.Context, userID string) (domain.UserHistory, error) {
	collection := uhr.database.Collection(uhr.collection)
	filter := bson.D{{Key: "_id", Value: userID}}

	var userHistory domain.UserHistory

	err := collection.FindOne(c, filter).Decode(&userHistory)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.UserHistory{
				UserID: userID,
				Items:  make([]domain.HistoryItem, 0),
			}, nil
		}
		return userHistory, err
	}

	return userHistory, nil
}
