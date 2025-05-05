package controller

import (
	"encoding/json"
	"main/domain"
	"main/internal"
	"net/http"
	"strconv"
	"strings"
)

type UserController struct {
	UserUseCase       domain.UserUseCase
	CollectionUseCase domain.CollectionUseCase
	HistoryUseCase    domain.HistoryUseCase
}

func (uc *UserController) Get(w http.ResponseWriter, r *http.Request) {
	authID := r.Context().Value("x-user-id").(string)

	queryParams := r.URL.Query()
	id := queryParams.Get("id")

	switch {
	case id == "": // получение своего профиля
		user, err := uc.UserUseCase.GetByID(r.Context(), authID)
		if err != nil {
			http.Error(w, jsonError(err.Error()), http.StatusNotFound)
			return
		}

		userInfo := domain.UserInfo{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			HasPicture:  user.HasPicture,
			Collections: user.Collections,
			Statistics:  user.Statistics,
			Favourite:   user.Favourite,
		}
		if user.Collections == nil {
			userInfo.Collections = make([]string, 0)
		}
		if user.Favourite == nil {
			userInfo.Favourite = make([]string, 0)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userInfo)
	case internal.ValidateUUID(id) != nil:
		http.Error(w, jsonError("Invalid id"), http.StatusBadRequest)
		return
	default: // получение чужого профиля
		user, err := uc.UserUseCase.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, jsonError(err.Error()), http.StatusNotFound)
			return
		}

		publicCollections := make([]string, 0)

		collections, err := uc.CollectionUseCase.SearchPublicByAuthor(r.Context(), id)
		if err != nil {
			http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
			return
		}
		for _, coll := range collections {
			publicCollections = append(publicCollections, coll.ID)
		}

		publicUserInfo := domain.PublicUserInfo{
			ID:                user.ID,
			Username:          user.Username,
			HasPicture:        user.HasPicture,
			PublicCollections: publicCollections,
			Statistics:        user.Statistics,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(publicUserInfo)
	}
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
	authID := r.Context().Value("x-user-id").(string)

	queryParams := r.URL.Query()
	id := queryParams.Get("id")

	if id == "" {
		id = authID // получим свою аватарку
	} else if internal.ValidateUUID(id) != nil {
		http.Error(w, jsonError("Invalid id"), http.StatusBadRequest)
		return
	}

	fileBytes, err := uc.UserUseCase.GetProfilePicture(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("User picture not found"), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}

//nolint:mnd // 5 MB
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
	//nolint:staticcheck // business logic
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

func (uc *UserController) GetHistory(w http.ResponseWriter, r *http.Request) {
	var result domain.UserHistoryArray
	var fromTime = 0
	var err error

	userID := r.Context().Value("x-user-id").(string)

	queryParams := r.URL.Query()
	fromTimeStr := queryParams.Get("from_time")

	if fromTimeStr != "" {
		fromTime, err = strconv.Atoi(fromTimeStr)

		if err != nil || fromTime < 0 {
			http.Error(w, jsonError("Invalid from_time"), http.StatusBadRequest)
			return
		}
	}

	userHistory, err := uc.HistoryUseCase.GetUserHistoryFromTime(r.Context(), userID, fromTime)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
		return
	}

	result.Count = len(userHistory)
	result.Items = userHistory

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
