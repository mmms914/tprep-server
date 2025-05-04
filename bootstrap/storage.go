package bootstrap

import (
	"context"
	"github.com/gookit/slog"
	"main/domain"
	"main/storage"
	"os"
	"time"
)

func NewStorage(env *Env) storage.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var endpoint string

	endpoint, exists := os.LookupEnv("MINIO_URI")
	if !exists {
		slog.Fatal("Cannot find MINIO_URI system variable")
	}

	accessKeyID := env.MinioRootUser
	secretAccessKey := env.MinioRootPassword
	useSSL := false

	options := storage.Options{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		UseSSL:          useSSL,
	}
	minioClient, err := storage.New(endpoint, options)
	if err != nil {
		slog.Fatal(err)
	}

	found, err := minioClient.BucketExists(ctx, domain.UserBucket)
	if err != nil {
		slog.Fatal(err)
	}
	if !found {
		err = minioClient.MakeBucket(ctx, domain.UserBucket)
		if err != nil {
			slog.Fatal("Can't create bucket:", err)
		}
		slog.Infof("Created bucket %q\n", domain.UserBucket)
	}

	found, err = minioClient.BucketExists(ctx, domain.CollectionBucket)
	if err != nil {
		slog.Fatal(err)
	}
	if !found {
		err = minioClient.MakeBucket(ctx, domain.CollectionBucket)
		if err != nil {
			slog.Fatal("Can't create bucket:", err)
		}
		slog.Infof("Created bucket %q\n", domain.CollectionBucket)
	}
	slog.Println("Connected to Minio")
	return minioClient
}
