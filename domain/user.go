package domain

import (
	"context"
	"io"
	"main/database"
)

const (
	UserCollection = "users"
	UserBucket     = "users"
)

type User struct {
	ID          string         `bson:"_id"         json:"id"`
	Username    string         `bson:"username"    json:"username"`
	Email       string         `bson:"email"       json:"email"`
	Password    string         `bson:"password"    json:"password"`
	HasPicture  bool           `bson:"has_picture" json:"has_picture"`
	Collections []string       `bson:"collections" json:"collections"`
	Favourite   []string       `bson:"favourite"   json:"favourite"`
	Statistics  UserStatistics `bson:"statistics"  json:"statistics"`
	Limits      UserLimits     `bson:"limits"      json:"limits"`
}

type UserInfo struct {
	ID          string         `bson:"_id"         json:"id"`
	Username    string         `bson:"username"    json:"username"`
	Email       string         `bson:"email"       json:"email"`
	HasPicture  bool           `bson:"has_picture" json:"has_picture"`
	Collections []string       `bson:"collections" json:"collections"`
	Statistics  UserStatistics `bson:"statistics"  json:"statistics"`
	Favourite   []string       `bson:"favourite"   json:"favourite"`
}

type PublicUserInfo struct {
	ID                string         `bson:"_id"         json:"id"`
	Username          string         `bson:"username"    json:"username"`
	HasPicture        bool           `bson:"has_picture" json:"has_picture"`
	PublicCollections []string       `bson:"collections" json:"collections"`
	Statistics        UserStatistics `bson:"statistics"  json:"statistics"`
}

type UserRepository interface {
	Create(c context.Context, user *User) (string, error)
	Update(c context.Context, filter interface{}, update interface{}) (database.UpdateResult, error)
	UpdateByID(c context.Context, userID string, update interface{}) (database.UpdateResult, error)
	DeleteByID(c context.Context, userID string) error
	GetByID(c context.Context, userID string) (User, error)
	GetByEmail(c context.Context, email string) (User, error)
	AddCollection(c context.Context, userID string, collectionID string, collectionType string) error
	DeleteCollection(c context.Context, userID string, collectionID string, collectionType string) error
}

type UserUseCase interface {
	PutByID(c context.Context, userID string, user *User) error
	GetByID(c context.Context, userID string) (User, error)
	DeleteByID(c context.Context, userID string) error
	GetProfilePicture(c context.Context, userID string) ([]byte, error)
	UploadProfilePicture(c context.Context, userID string, picture io.Reader, size int64) error
	RemoveProfilePicture(c context.Context, userID string) error
}

//nolint:iface // business logic
type UserStorage interface {
	GetObject(c context.Context, objectName string) ([]byte, error)
	PutObject(c context.Context, objectName string, reader io.Reader, objectSize int64) error
	RemoveObject(c context.Context, objectName string) error
}
