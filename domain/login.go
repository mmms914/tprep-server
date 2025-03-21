package domain

import (
	"context"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginUseCase interface {
	GetUserByEmail(c context.Context, email string) (User, error)
	CreateAccessToken(user *User, secret string, expiry int) (accessToken string, expTime time.Time, err error)
	CreateRefreshToken(user *User, secret string, expiry int) (refreshToken string, expTime time.Time, err error)
}
