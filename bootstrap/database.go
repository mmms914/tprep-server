package bootstrap

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
)

func NewMongoDatabase(env *Env) *mongo.Client {
	mongodbURI := env.MongoURI
	client, err := mongo.Connect(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatal(err)
	}
	return client
}
