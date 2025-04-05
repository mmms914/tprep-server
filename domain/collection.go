package domain

import (
	"context"
	"io"
	"main/database"
)

const (
	CollectionCollection = "collections"
	CollectionBucket     = "collections"
)

type Collection struct {
	ID        string `bson:"_id" json:"id"`
	Name      string `bson:"name" json:"name"`
	IsPublic  bool   `bson:"is_public" json:"is_public"`
	Cards     []Card `bson:"cards" json:"cards"`
	MaxId     int    `bson:"max_id" json:"max_id"`
	Author    string `bson:"author" json:"author"`
	Likes     int    `bson:"likes" json:"likes"`
	Trainings int    `bson:"trainings" json:"trainings"`
}

type CollectionInfo struct {
	ID        string `bson:"_id" json:"id"`
	Name      string `bson:"name" json:"name"`
	IsPublic  bool   `bson:"is_public" json:"is_public"`
	Cards     []Card `bson:"cards" json:"cards"`
	Author    string `bson:"author" json:"author"`
	Likes     int    `bson:"likes" json:"likes"`
	Trainings int    `bson:"trainings" json:"trainings"`
}

type CollectionPreview struct {
	ID         string `bson:"_id" json:"id"`
	Name       string `bson:"name" json:"name"`
	IsPublic   bool   `bson:"is_public" json:"is_public"`
	CardsCount int    `json:"cards_count"`
	Likes      int    `bson:"likes" json:"likes"`
	Trainings  int    `bson:"trainings" json:"trainings"`
}

type CollectionPreviewArray struct {
	Count int                 `json:"count"`
	Items []CollectionPreview `json:"items"`
}

type CollectionRepository interface {
	Create(c context.Context, collection *Collection) (string, error)
	Update(c context.Context, filter interface{}, update interface{}) (database.UpdateResult, error)
	UpdateByID(c context.Context, collectionID string, update interface{}) (database.UpdateResult, error)
	DeleteByID(c context.Context, collectionID string) error
	GetByID(c context.Context, collectionID string) (Collection, error)
	GetByFilter(c context.Context, filter interface{}, opts database.FindOptions) ([]Collection, error)
}

type CollectionUseCase interface {
	Create(c context.Context, collection *Collection) (string, error)
	PutByID(c context.Context, collectionID string, collection *Collection) error
	DeleteByID(c context.Context, collectionID string) error
	GetByID(c context.Context, collectionID string) (Collection, error)
	AddLike(c context.Context, collectionID string) (*Collection, error)
	RemoveLike(c context.Context, collectionID string) (*Collection, error)
	SearchPublic(c context.Context, text string, count int, offset int, sortBy string, category string, userID string) ([]Collection, error)
	SearchPublicByAuthor(c context.Context, author string) ([]Collection, error)
	AddCard(c context.Context, collectionID string, card *Card) (Card, error)
	DeleteCard(c context.Context, collectionID string, cardLocalID int) error
	UpdateCard(c context.Context, collectionID string, card *Card) error
	GetCardPhoto(c context.Context, objectName string) ([]byte, error)
	UploadCardPhoto(c context.Context, userID string, collectionID string, cardID int, picture io.Reader, size int64) (string, error)
	RemoveCardPicture(c context.Context, userID string, collectionID string, cardID int, objectName string) error
}

type CollectionStorage interface {
	GetObject(c context.Context, objectName string) ([]byte, error)
	PutObject(c context.Context, objectName string, reader io.Reader, objectSize int64) error
	RemoveObject(c context.Context, objectName string) error
}
