package usecase

import (
	"context"
	"main/domain"
	"main/internal"
	"time"
)

type signupUseCase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewSignupUseCase(userRepository domain.UserRepository, timeout time.Duration) domain.SignupUseCase {
	return &signupUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (su *signupUseCase) Create(c context.Context, user *domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	user.Collections = make([]string, 0)
	user.Favourite = make([]string, 0)
	return su.userRepository.Create(ctx, user)
}

func (su *signupUseCase) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.GetByEmail(ctx, email)
}

func (su *signupUseCase) CreateAccessToken(user *domain.User, secret string, expiry int) (string, time.Time, error) {
	accessToken, exp, err := internal.CreateAccessToken(user, secret, expiry)
	return accessToken, exp, err
}

func (su *signupUseCase) CreateRefreshToken(user *domain.User, secret string, expiry int) (string, time.Time, error) {
	refreshToken, exp, err := internal.CreateRefreshToken(user, secret, expiry)
	return refreshToken, exp, err
}
