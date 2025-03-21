package domain

import (
	"context"
	"time"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenUseCase interface {
	GetUserByID(c context.Context, id string) (User, error)
	CreateAccessToken(user *User, secret string, expiry int) (accessToken string, exp time.Time, err error)
	CreateRefreshToken(user *User, secret string, expiry int) (refreshToken string, exp time.Time, err error)
	ExtractIDFromToken(requestToken string, secret string) (string, error)
}
