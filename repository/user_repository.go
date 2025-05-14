package repository

import (
	"context"
	"errors"
	"main/database"
	"main/domain"
	"main/internal"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type userRepository struct {
	database   database.Database
	collection string
}

func NewUserRepository(db database.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

func (ur *userRepository) Create(c context.Context, user *domain.User) (string, error) {
	collection := ur.database.Collection(ur.collection)
	user.ID = internal.GenerateUUID()
	id, err := collection.InsertOne(c, user)
	return id, err
}

func (ur *userRepository) UpdateByID(
	c context.Context,
	userID string,
	update interface{},
) (database.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: userID}}
	return ur.Update(c, filter, update)
}

func (ur *userRepository) Update(
	c context.Context,
	filter interface{},
	update interface{},
) (database.UpdateResult, error) {
	collection := ur.database.Collection(ur.collection)

	return collection.UpdateOne(c, filter, update)
}

func (ur *userRepository) DeleteByID(c context.Context, userID string) error {
	collection := ur.database.Collection(ur.collection)
	filter := bson.D{{Key: "_id", Value: userID}}
	_, err := collection.DeleteOne(c, filter)
	return err
}

func (ur *userRepository) GetByID(c context.Context, userID string) (domain.User, error) {
	var user domain.User
	collection := ur.database.Collection(ur.collection)
	filter := bson.D{{Key: "_id", Value: userID}}
	err := collection.FindOne(c, filter).Decode(&user)
	return user, err
}

func (ur *userRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	var user domain.User
	collection := ur.database.Collection(ur.collection)
	filter := bson.D{{Key: "email", Value: email}}
	err := collection.FindOne(c, filter).Decode(&user)
	return user, err
}

func (ur *userRepository) AddCollection(
	c context.Context,
	userID string,
	collectionID string,
	collectionType string,
) error {
	var update bson.D
	if collectionType == "collections" {
		update = bson.D{
			{Key: "$push", Value: bson.D{
				{Key: "collections", Value: collectionID},
			}},
		}
	}
	if collectionType == "favourite" {
		update = bson.D{
			{Key: "$push", Value: bson.D{
				{Key: "favourite", Value: collectionID},
			}},
		}
	}

	res, err := ur.UpdateByID(c, userID, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (ur *userRepository) DeleteCollection(
	c context.Context,
	userID string,
	collectionID string,
	collectionType string,
) error {
	var update bson.D
	if collectionType == "collections" {
		update = bson.D{
			{Key: "$pull", Value: bson.D{
				{Key: "collections", Value: collectionID},
			}},
		}
	}
	if collectionType == "favourite" {
		update = bson.D{
			{Key: "$pull", Value: bson.D{
				{Key: "favourite", Value: collectionID},
			}},
		}
	}
	res, err := ur.UpdateByID(c, userID, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
