package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client interface {
	GetObject(ctx context.Context, bucketName string, objectName string) ([]byte, error)
	PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64) error
	RemoveObject(ctx context.Context, bucketName string, objectName string) error
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	MakeBucket(ctx context.Context, bucketName string) error
}

type storageClient struct {
	cl *minio.Client
}

type Options struct {
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

const jpegForm = ".jpeg"

func New(endpoint string, options Options) (Client, error) {
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(options.AccessKeyID, options.SecretAccessKey, ""),
		Secure: options.UseSSL,
	}
	cl, err := minio.New(endpoint, opts)
	return &storageClient{cl: cl}, err
}

func (sc *storageClient) GetObject(ctx context.Context, bucketName string, objectName string) ([]byte, error) {
	objectName += jpegForm
	obj, err := sc.cl.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	fileBytes, err := io.ReadAll(obj)
	return fileBytes, err
}

func (sc *storageClient) PutObject(
	ctx context.Context,
	bucketName string,
	objectName string,
	reader io.Reader,
	objectSize int64,
) error {
	objectName += jpegForm
	_, err := sc.cl.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{})
	return err
}

func (sc *storageClient) RemoveObject(ctx context.Context, bucketName string, objectName string) error {
	objectName += jpegForm
	return sc.cl.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (sc *storageClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return sc.cl.BucketExists(ctx, bucketName)
}

func (sc *storageClient) MakeBucket(ctx context.Context, bucketName string) error {
	return sc.cl.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}
