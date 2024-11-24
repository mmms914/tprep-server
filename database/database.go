package database

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"main/models"
	"slices"
)

var database *mongo.Database
var globalValues models.GlobalValues
var collections, collectionsGlobal, globals *mongo.Collection

func GlobalValues() models.GlobalValues {
	err := globals.FindOne(context.TODO(), bson.D{}).Decode(&globalValues)
	if err != nil {
		log.Fatal("cannot get global values")
	}

	return globalValues
}

func SetGlobalValues(gv models.GlobalValues) {
	globalValues = gv
	globals.FindOneAndReplace(context.TODO(), bson.D{}, gv)
}

func InitDatabase(db *mongo.Database) (*mongo.Collection, *mongo.Collection) {
	database = db
	names, err := database.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	globals = database.Collection("globals")
	if !slices.Contains(names, "globals") {
		_, err := globals.InsertOne(context.TODO(), globalValues)
		if err != nil {
			log.Fatal(err)
		}
	}

	collections = database.Collection("collections")
	collectionsGlobal = database.Collection("collections_global")

	return collections, collectionsGlobal
}
