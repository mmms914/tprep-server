package usecase

import (
	"context"
	"main/domain"
	"main/internal"
	"time"
)

type refreshTokenUseCase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewRefreshTokenUseCase(userRepository domain.UserRepository, timeout time.Duration) domain.RefreshTokenUseCase {
	return &refreshTokenUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (rtu *refreshTokenUseCase) GetUserByID(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, rtu.contextTimeout)
	defer cancel()
	return rtu.userRepository.GetByID(ctx, email)
}

func (rtu *refreshTokenUseCase) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, exp time.Time, err error) {
	return internal.CreateAccessToken(user, secret, expiry)
}

func (rtu *refreshTokenUseCase) CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, exp time.Time, err error) {
	return internal.CreateRefreshToken(user, secret, expiry)
}

func (rtu *refreshTokenUseCase) ExtractIDFromToken(requestToken string, secret string) (string, error) {
	return internal.ExtractIDFromToken(requestToken, secret)
}
