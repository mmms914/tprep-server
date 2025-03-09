package controller

import (
	"encoding/json"
	"main/domain"
	"net/http"
	"strings"
)

type UserController struct {
	UserUseCase domain.UserUseCase
}

func (uc *UserController) Get(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("x-user-id").(string)

	user, err := uc.UserUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
		return
	}

	userInfo := domain.UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Collections: user.Collections,
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
		http.Error(w, jsonError("Invalid data"), http.StatusBadRequest)
		return
	}

	id := r.Context().Value("x-user-id").(string)
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

func (uc *UserController) GetProfilePicture(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("x-user-id").(string)
	fileBytes, err := uc.UserUseCase.GetProfilePicture(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("User picture not found"), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}

func (uc *UserController) UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, jsonError("Error with file or its max size"), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if handler.Size > (5 << 20) {
		http.Error(w, jsonError("Image's size should be less than 5 MB"), http.StatusBadRequest)
		return
	}

	if !(strings.HasSuffix(handler.Filename, ".jpg") || strings.HasSuffix(handler.Filename, ".jpeg")) {
		http.Error(w, jsonError("Image's extension should be JPG/JPEG"), http.StatusBadRequest)
		return
	}

	id := r.Context().Value("x-user-id").(string)
	err = uc.UserUseCase.UploadProfilePicture(r.Context(), id, file, handler.Size)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Picture uploaded successfully",
	})
}

func (uc *UserController) RemoveProfilePicture(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("x-user-id").(string)
	err := uc.UserUseCase.RemoveProfilePicture(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Profile picture deleted",
	})
}
