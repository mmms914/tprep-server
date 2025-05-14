package controller

import (
	"encoding/json"
	"main/bootstrap"
	"main/domain"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	LoginUseCase domain.LoginUseCase
	Env          *bootstrap.Env
}

func (lc *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	var request domain.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	user, err := lc.LoginUseCase.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		http.Error(w, jsonError("User with this email not found"), http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		http.Error(w, jsonError("Invalid credentials"), http.StatusBadRequest)
		return
	}

	accessToken, expAccess, err := lc.LoginUseCase.CreateAccessToken(
		&user,
		lc.Env.AccessTokenSecret,
		lc.Env.AccessTokenExpiryHour,
	)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	refreshToken, expRefresh, err := lc.LoginUseCase.CreateRefreshToken(
		&user,
		lc.Env.RefreshTokenSecret,
		lc.Env.RefreshTokenExpiryHour,
	)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	loginResponse := domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Access-Expires-After", expAccess.UTC().String())
	w.Header().Set("X-Refresh-Expires-After", expRefresh.UTC().String())
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse)
}
