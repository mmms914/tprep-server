package bootstrap

import (
	"context"
	"main/database"
	"main/repository"
	"os"
	"time"

	"github.com/gookit/slog"
)

func NewMongoDatabase(_ *Env) database.Client {
	//nolint:mnd // 10 sec
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongodbURI, exists := os.LookupEnv("MONGODB_URI")
	if !exists {
		slog.Fatal("Cannot find MONGODB_URI system variable")
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
