package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"main/domain"
	"net/http"
)

type UserController struct {
	UserUseCase domain.UserUseCase
}

func (uc *UserController) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := uc.UserUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
		return
	}

	userInfo := domain.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userInfo)
}

func (uc *UserController) Update(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Email == "" {
		http.Error(w, jsonError("Invalid username"), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	err = uc.UserUseCase.PutByID(r.Context(), id, &user)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "User updated",
	})
}
