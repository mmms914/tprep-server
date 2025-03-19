package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	filter := bson.M{"user_id": userID}
	if collection.FindOne(c, filter).Decode(&userHistory) != nil {
		userHistory = domain.UserHistory{
			UserID: userID,
			Items:  make([]domain.HistoryItem, 0),
		}
		_, err := collection.InsertOne(c, userHistory)
		if err != nil {
			return err
		}
	}
	return nil
}

func (uhr *userHistoryRepository) UpdateByID(c context.Context, userID string, item domain.HistoryItem) error {
	err := uhr.CreateIfNotExists(c, userID)
	if err != nil {
		return err
	}

	collection := uhr.database.Collection(uhr.collection)
	filter := bson.D{{Key: "user_id", Value: userID}}
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
	filter := bson.D{{Key: "user_id", Value: userID}}

	var items []domain.HistoryItem

	cursor, err := collection.Find(c, filter)
	if err != nil {
		return domain.UserHistory{}, err
	}

	err = cursor.All(c, &items)
	if err != nil {
		return domain.UserHistory{}, err
	}

	userHistory := domain.UserHistory{
		UserID: userID,
		Items:  items,
	}
	return userHistory, nil
}
