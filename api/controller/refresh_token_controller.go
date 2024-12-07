package controller

import (
	"encoding/json"
	"main/bootstrap"
	"main/domain"
	"net/http"
)

type RefreshTokenController struct {
	RefreshTokenUseCase domain.RefreshTokenUseCase
	Env                 *bootstrap.Env
}

func (rtc *RefreshTokenController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var request domain.RefreshTokenRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	id, err := rtc.RefreshTokenUseCase.ExtractIDFromToken(request.RefreshToken, rtc.Env.RefreshTokenSecret)
	if err != nil {
		http.Error(w, jsonError("Invalid refresh token"), http.StatusUnauthorized)
		return
	}

	user, err := rtc.RefreshTokenUseCase.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("User not found"), http.StatusUnauthorized)
		return
	}

	accessToken, err := rtc.RefreshTokenUseCase.CreateAccessToken(&user, rtc.Env.AccessTokenSecret, rtc.Env.AccessTokenExpiryHour)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	refreshToken, err := rtc.RefreshTokenUseCase.CreateRefreshToken(&user, rtc.Env.RefreshTokenSecret, rtc.Env.RefreshTokenExpiryHour)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	refreshTokenResponse := domain.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(refreshTokenResponse)
}
