package controller

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"main/bootstrap"
	"main/domain"
	"net/http"
)

type SignupController struct {
	SignupUseCase domain.SignupUseCase
	Env           *bootstrap.Env
}

func (sc *SignupController) Signup(w http.ResponseWriter, r *http.Request) {
	var request domain.SignupRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if request.Username == "" || request.Email == "" || request.Password == "" {
		http.Error(w, jsonError("Invalid data"), http.StatusBadRequest)
		return
	}

	_, err = sc.SignupUseCase.GetUserByEmail(r.Context(), request.Email)
	if err == nil {
		http.Error(w, jsonError("User with this email already exists"), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	user := &domain.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	userID, err := sc.SignupUseCase.Create(r.Context(), user)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	accessToken, expAccess, err := sc.SignupUseCase.CreateAccessToken(user, sc.Env.AccessTokenSecret, sc.Env.AccessTokenExpiryHour)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	refreshToken, expRefresh, err := sc.SignupUseCase.CreateRefreshToken(user, sc.Env.RefreshTokenSecret, sc.Env.RefreshTokenExpiryHour)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	signupResponse := domain.SignupResponse{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Access-Expires-After", expAccess.UTC().String())
	w.Header().Set("X-Refresh-Expires-After", expRefresh.UTC().String())
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(signupResponse)
}
