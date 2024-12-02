package domain

import "context"

type Collection struct {
	ID       string `bson:"id" json:"id"`
	Name     string `bson:"name" json:"name"`
	IsPublic bool   `bson:"is_public" json:"is_public"`
	Cards    []Card `bson:"cards" json:"cards"`
	MaxId    int    `bson:"max_id" json:"max_id"`
}

type CollectionRepository interface {
	Create(c context.Context, collection *Collection) error
	Update(c context.Context, filter interface{}, update interface{}) error
	UpdateByID(c context.Context, collectionID string, update interface{}) error
	DeleteByID(c context.Context, collectionID string) error
	GetByID(c context.Context, collectionID string) (*Collection, error)
}

type CollectionUseCase interface {
	Create(c context.Context, collection *Collection) error
	UpdateNameByID(c context.Context, collectionID string, collectionName string) error
	UpdateTypeByID(c context.Context, collectionID string, collectionType bool) error
	DeleteByID(c context.Context, collectionID string) error
	GetByID(c context.Context, collectionID string) (*Collection, error)

	AddCard(c context.Context, collectionID string, card *Card) error
	DeleteCard(c context.Context, collectionID string, cardLocalID int) error
	UpdateCard(c context.Context, collectionID string, card *Card) error
}
