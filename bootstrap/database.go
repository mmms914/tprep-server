package bootstrap

import (
	"context"
	"github.com/gookit/slog"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoDatabase(env *Env) *mongo.Client {
	var mongodbURI string
	if env.AppEnv == "local" {
		mongodbURI = env.LocalMongoURI
	} else if env.AppEnv == "docker" {
		mongodbURI = env.DockerMongoURI
	}

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
