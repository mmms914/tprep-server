package domain

import (
	"context"
	"time"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	UserID       string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignupUseCase interface {
	Create(c context.Context, user *User) (string, error)
	GetUserByEmail(c context.Context, email string) (User, error)
	CreateAccessToken(user *User, secret string, expiry int) (accessToken string, exp time.Time, err error)
	CreateRefreshToken(user *User, secret string, expiry int) (refreshToken string, exp time.Time, err error)
}
