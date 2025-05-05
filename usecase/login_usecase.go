package usecase

import (
	"context"
	"main/domain"
	"main/internal"
	"time"
)

type loginUseCase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewLoginUseCase(userRepository domain.UserRepository, timeout time.Duration) domain.LoginUseCase {
	return &loginUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (lu *loginUseCase) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.GetByEmail(ctx, email)
}

func (lu *loginUseCase) CreateAccessToken(user *domain.User, secret string, expiry int) (string, time.Time, error) {
	accessToken, exp, err := internal.CreateAccessToken(user, secret, expiry)
	return accessToken, exp, err
}

func (lu *loginUseCase) CreateRefreshToken(user *domain.User, secret string, expiry int) (string, time.Time, error) {
	refreshToken, exp, err := internal.CreateRefreshToken(user, secret, expiry)
	return refreshToken, exp, err
}
