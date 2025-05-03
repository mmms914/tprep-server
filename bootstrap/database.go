package bootstrap

import (
	"context"
	"github.com/gookit/slog"
	"main/database"
	"main/repository"
	"time"
)

func NewMongoDatabase(env *Env) database.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var mongodbURI string
	if env.AppEnv == "local" {
		mongodbURI = env.LocalMongoURI
	} else if env.AppEnv == "docker" {
		mongodbURI = env.DockerMongoURI
	}

	client, err := database.NewClient(mongodbURI)
	if err != nil {
		slog.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		slog.Fatal(err)
	}

	repository.SetClient(client)

	slog.Println("Successfully connected to MongoDB!")
	return client
}

func CloseMongoDBConnection(client database.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		slog.Fatal(err)
	}

	slog.Println("Connection to MongoDB closed.")
}
