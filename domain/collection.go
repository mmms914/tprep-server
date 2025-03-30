package domain

import (
	"context"
	"main/database"
)

const (
	CollectionCollection = "collections"
)

type Collection struct {
	ID       string `bson:"_id" json:"id"`
	Name     string `bson:"name" json:"name"`
	IsPublic bool   `bson:"is_public" json:"is_public"`
	Cards    []Card `bson:"cards" json:"cards"`
	MaxId    int    `bson:"max_id" json:"max_id"`
	Author   string `bson:"author" json:"author"`
	Likes    int    `bson:"likes" json:"likes"`
}

type CollectionInfo struct {
	ID       string `bson:"_id" json:"id"`
	Name     string `bson:"name" json:"name"`
	IsPublic bool   `bson:"is_public" json:"is_public"`
	Cards    []Card `bson:"cards" json:"cards"`
	Author   string `bson:"author" json:"author"`
	Likes    int    `bson:"likes" json:"likes"`
}

type CollectionPreview struct {
	ID         string `bson:"_id" json:"id"`
	Name       string `bson:"name" json:"name"`
	IsPublic   bool   `bson:"is_public" json:"is_public"`
	CardsCount int    `json:"cards_count"`
	Likes      int    `bson:"likes" json:"likes"`
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
	SearchPublic(c context.Context, text string, count int, offset int) ([]Collection, error)
	SearchPublicByAuthor(c context.Context, author string) ([]Collection, error)
	AddCard(c context.Context, collectionID string, card *Card) (Card, error)
	DeleteCard(c context.Context, collectionID string, cardLocalID int) error
	UpdateCard(c context.Context, collectionID string, card *Card) error
}
