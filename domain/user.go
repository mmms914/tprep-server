package domain

import "context"

const (
	UserCollection = "users"
)

type User struct {
	ID          string   `bson:"_id" json:"id"`
	Username    string   `bson:"username" json:"username"`
	Email       string   `bson:"email" json:"email"`
	Password    string   `bson:"password" json:"password"`
	Collections []string `bson:"collections" json:"collections"`
}

type UserInfo struct {
	ID          string   `bson:"_id" json:"id"`
	Username    string   `bson:"username" json:"username"`
	Email       string   `bson:"email" json:"email"`
	Collections []string `bson:"collections" json:"collections"`
}

type UserRepository interface {
	Create(c context.Context, user *User) (string, error)
	Update(c context.Context, filter interface{}, update interface{}) error
	UpdateByID(c context.Context, userID string, update interface{}) error
	DeleteByID(c context.Context, userID string) error
	GetByID(c context.Context, userID string) (User, error)
	GetByEmail(c context.Context, email string) (User, error)
}

type UserUseCase interface {
	PutByID(c context.Context, userID string, user *User) error
	GetByID(c context.Context, userID string) (User, error)
	DeleteByID(c context.Context, userID string) error
	AddCollection(c context.Context, userID string, collectionID string) error
	DeleteCollection(c context.Context, userID string, collectionID string) error
}
