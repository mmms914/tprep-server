package storage

import (
	"context"
	"io"
	"main/domain"
)

type userStorage struct {
	storage Client
	bucket  string
}

func (us *userStorage) GetObject(c context.Context, objectName string) ([]byte, error) {
	return us.storage.GetObject(c, us.bucket, objectName)
}

func (us *userStorage) PutObject(c context.Context, objectName string, reader io.Reader, objectSize int64) error {
	return us.storage.PutObject(c, us.bucket, objectName, reader, objectSize)
}

func (us *userStorage) RemoveObject(c context.Context, objectName string) error {
	return us.storage.RemoveObject(c, us.bucket, objectName)
}

func NewUserStorage(s Client, bucket string) domain.UserStorage {
	return &userStorage{
		storage: s,
		bucket:  bucket,
	}
}
