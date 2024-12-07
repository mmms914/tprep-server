package usecase

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"main/domain"
	"time"
)

type userUseCase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUseCase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (uu *userUseCase) PutByID(c context.Context, userID string, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	update := bson.D{
		{"$set", bson.D{
			{"username", user.Username},
			{"email", user.Email},
		}},
	}
	return uu.userRepository.UpdateByID(ctx, userID, update)
}

func (uu *userUseCase) DeleteByID(c context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.DeleteByID(ctx, userID)
}

func (uu *userUseCase) GetByID(c context.Context, userID string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.GetByID(ctx, userID)
}

func (uu *userUseCase) AddCollection(c context.Context, userID string, collectionID string) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	update := bson.D{
		{"$push", bson.D{
			{"collections", collectionID},
		}},
	}
	return uu.userRepository.UpdateByID(ctx, userID, update)
}

func (uu *userUseCase) DeleteCollection(c context.Context, userID string, collectionID string) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	update := bson.D{
		{"$pull", bson.D{
			{"collections", collectionID},
		}},
	}
	return uu.userRepository.UpdateByID(ctx, userID, update)
}
