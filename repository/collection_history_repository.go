package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"main/database"
	"main/domain"
)

type collectionHistoryRepository struct {
	database   database.Database
	collection string
}

func NewCollectionHistoryRepository(db database.Database, collection string) domain.CollectionHistoryRepository {
	return &collectionHistoryRepository{
		database:   db,
		collection: collection,
	}
}

func (chr *collectionHistoryRepository) CreateIfNotExists(c context.Context, collectionID string) error {
	collection := chr.database.Collection(chr.collection)

	var collectionHistory domain.CollectionHistory
	filter := bson.M{"collection_id": collectionID}

	if collection.FindOne(c, filter).Decode(&collectionHistory) != nil {
		collectionHistory = domain.CollectionHistory{
			CollectionID: collectionID,
			Items:        make([]domain.SmallHistoryItem, 0),
		}
		_, err := collection.InsertOne(c, collectionHistory)
		if err != nil {
			return err
		}
	}
	return nil
}

func (chr *collectionHistoryRepository) UpdateByID(c context.Context, collectionID string, item domain.SmallHistoryItem) error {
	err := chr.CreateIfNotExists(c, collectionID)
	if err != nil {
		return err
	}

	collection := chr.database.Collection(chr.collection)
	filter := bson.D{{Key: "collection_id", Value: collectionID}}
	update := bson.D{
		{"$push", bson.D{
			{"items", item},
		}},
	}
	_, err = collection.UpdateOne(c, filter, update)

	return err
}

func (chr *collectionHistoryRepository) GetByID(c context.Context, collectionID string) (domain.CollectionHistory, error) {
	collection := chr.database.Collection(chr.collection)
	filter := bson.D{{Key: "collection_id", Value: collectionID}}

	var items []domain.SmallHistoryItem

	cursor, err := collection.Find(c, filter)
	if err != nil {
		return domain.CollectionHistory{}, err
	}

	err = cursor.All(c, &items)
	if err != nil {
		return domain.CollectionHistory{}, err
	}

	collectionHistory := domain.CollectionHistory{
		CollectionID: collectionID,
		Items:        items,
	}
	return collectionHistory, nil
}
