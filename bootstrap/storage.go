package bootstrap

import (
	"context"
	"github.com/gookit/slog"
	"main/domain"
	"main/storage"
	"time"
)

func NewStorage(env *Env) storage.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var endpoint string

	if env.AppEnv == "local" {
		endpoint = env.LocalMinioURI
	} else if env.AppEnv == "docker" {
		endpoint = env.DockerMinioURI
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
	slog.Println("Connected to Minio")
	return minioClient
}
