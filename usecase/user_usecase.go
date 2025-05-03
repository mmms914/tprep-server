package usecase

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"io"
	"main/domain"
	"time"
)

type userUseCase struct {
	userRepository domain.UserRepository
	userStorage    domain.UserStorage
	contextTimeout time.Duration
}

func NewUserUseCase(userRepository domain.UserRepository, userStorage domain.UserStorage, timeout time.Duration) domain.UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
		userStorage:    userStorage,
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
	_, err := uu.userRepository.UpdateByID(ctx, userID, update)
	return err
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

func (uu *userUseCase) GetProfilePicture(c context.Context, userID string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	return uu.userStorage.GetObject(ctx, userID)
}

func (uu *userUseCase) UploadProfilePicture(c context.Context, userID string, picture io.Reader, size int64) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	update := bson.D{
		{"$set", bson.D{
			{"has_picture", true},
		}},
	}

	_, err := uu.userRepository.UpdateByID(ctx, userID, update)
	if err != nil {
		return err
	}

	return uu.userStorage.PutObject(ctx, userID, picture, size)
}

func (uu *userUseCase) RemoveProfilePicture(c context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	update := bson.D{
		{"$set", bson.D{
			{"has_picture", false},
		}},
	}

	_, err := uu.userRepository.UpdateByID(ctx, userID, update)
	if err != nil {
		return err
	}

	return uu.userStorage.RemoveObject(ctx, userID)
}
