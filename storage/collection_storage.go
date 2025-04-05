package storage

import (
	"context"
	"io"
	"main/domain"
)

type collectionStorage struct {
	storage Client
	bucket  string
}

func (cs *collectionStorage) GetObject(c context.Context, objectName string) ([]byte, error) {
	return cs.storage.GetObject(c, cs.bucket, objectName)
}

func (cs *collectionStorage) PutObject(c context.Context, objectName string, reader io.Reader, objectSize int64) error {
	return cs.storage.PutObject(c, cs.bucket, objectName, reader, objectSize)
}

func (cs *collectionStorage) RemoveObject(c context.Context, objectName string) error {
	return cs.storage.RemoveObject(c, cs.bucket, objectName)
}

func NewCollectionStorage(s Client, bucket string) domain.CollectionStorage {
	return &collectionStorage{
		storage: s,
		bucket:  bucket,
	}
}
