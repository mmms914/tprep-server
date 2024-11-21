package bootstrap

import (
	"context"
	"github.com/gookit/slog"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoDatabase(env *Env) *mongo.Client {
	mongodbURI := env.MongoURI
	client, err := mongo.Connect(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		slog.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		slog.FatalErr(err)
	}
	slog.Println("Successfully connected to MongoDB!")
	return client
}
