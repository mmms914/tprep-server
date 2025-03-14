package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"main/domain"
	"net/http"
	"strconv"
)

type CollectionController struct {
	CollectionUseCase domain.CollectionUseCase
	UserUseCase       domain.UserUseCase
}

func (cc *CollectionController) Create(w http.ResponseWriter, r *http.Request) {
	var collection domain.Collection

	err := json.NewDecoder(r.Body).Decode(&collection)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if collection.Name == "" {
		http.Error(w, jsonError("Invalid name"), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("x-user-id").(string)
	collection.Author = userID

	id, err := cc.CollectionUseCase.Create(r.Context(), &collection)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	collection, err = cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	err = cc.UserUseCase.AddCollection(r.Context(), userID, id)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	collectionInfo := domain.CollectionInfo{
		ID:       collection.ID,
		Name:     collection.Name,
		IsPublic: collection.IsPublic,
		Cards:    collection.Cards,
		Author:   collection.Author,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(collectionInfo)
}

func (cc *CollectionController) Update(w http.ResponseWriter, r *http.Request) {
	var collection domain.Collection
	userID := r.Context().Value("x-user-id").(string)

	err := json.NewDecoder(r.Body).Decode(&collection)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if collection.Name == "" {
		http.Error(w, jsonError("Invalid name"), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	coll, err := cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if coll.Author != userID {
		http.Error(w, jsonError("You are not the owner of this collection"), http.StatusForbidden)
		return
	}

	err = cc.CollectionUseCase.PutByID(r.Context(), id, &collection)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Collection updated",
	})
}

func (cc *CollectionController) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("x-user-id").(string)

	id := chi.URLParam(r, "id")
	collection, err := cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if userID != collection.Author && collection.IsPublic == false {
		http.Error(w, jsonError("You are not the owner of this collection"), http.StatusForbidden)
		return
	}

	collectionInfo := domain.CollectionInfo{
		ID:       collection.ID,
		Name:     collection.Name,
		IsPublic: collection.IsPublic,
		Cards:    collection.Cards,
		Author:   collection.Author,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(collectionInfo)
}

func (cc *CollectionController) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("x-user-id").(string)

	id := chi.URLParam(r, "id")

	coll, err := cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if coll.Author != userID {
		http.Error(w, jsonError("You are not the owner of this collection"), http.StatusForbidden)
		return
	}

	err = cc.CollectionUseCase.DeleteByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
		return
	}

	err = cc.UserUseCase.DeleteCollection(r.Context(), userID, id)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Collection deleted successfully",
	})
}

func (cc *CollectionController) Search(w http.ResponseWriter, r *http.Request) {
	var collections []domain.Collection
	var result domain.CollectionPreviewArray
	var err error

	queryParams := r.URL.Query()
	name := queryParams.Get("name")
	count, err := strconv.Atoi(queryParams.Get("count"))
	if err != nil || count < 1 || count > 100 {
		http.Error(w, jsonError("Invalid count"), http.StatusBadRequest)
		return
	}
	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil || offset < 0 {
		http.Error(w, jsonError("Invalid offset"), http.StatusBadRequest)
		return
	}

	collections, err = cc.CollectionUseCase.SearchPublic(r.Context(), name, count, offset)

	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	if len(collections) == 0 {
		http.Error(w, "Сouldn't find anything", http.StatusNotFound)
		return
	}

	cnt := 0
	for _, c := range collections {
		result.Items = append(result.Items, domain.CollectionPreview{
			ID:         c.ID,
			Name:       c.Name,
			IsPublic:   c.IsPublic,
			CardsCount: len(c.Cards),
		})
		cnt++
	}
	result.Count = cnt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (cc *CollectionController) CreateCard(w http.ResponseWriter, r *http.Request) {
	var card domain.Card
	userID := r.Context().Value("x-user-id").(string)

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}
	if card.Question == "" || card.Answer == "" {
		http.Error(w, jsonError("Invalid body data"), http.StatusBadRequest)
		return
	}
	collectionID := chi.URLParam(r, "id")

	coll, err := cc.CollectionUseCase.GetByID(r.Context(), collectionID)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if coll.Author != userID {
		http.Error(w, jsonError("You are not the owner of this collection"), http.StatusForbidden)
		return
	}

	card, err = cc.CollectionUseCase.AddCard(r.Context(), collectionID, &card)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

func (cc *CollectionController) UpdateCard(w http.ResponseWriter, r *http.Request) {
	var card domain.Card
	userID := r.Context().Value("x-user-id").(string)

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if card.Question == "" || card.Answer == "" {
		http.Error(w, jsonError("Invalid body data"), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	card.LocalID, err = strconv.Atoi(chi.URLParam(r, "cardID"))
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	coll, err := cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if coll.Author != userID {
		http.Error(w, jsonError("You are not the owner of this collection"), http.StatusForbidden)
		return
	}

	err = cc.CollectionUseCase.UpdateCard(r.Context(), id, &card)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Card updated",
	})
}

func (cc *CollectionController) DeleteCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("x-user-id").(string)

	collectionID := chi.URLParam(r, "id")
	cardID, err := strconv.Atoi(chi.URLParam(r, "cardID"))
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	coll, err := cc.CollectionUseCase.GetByID(r.Context(), collectionID)
	if err != nil {
		http.Error(w, jsonError("There is no collection with this ID"), http.StatusNotFound)
		return
	}

	if coll.Author != userID {
		http.Error(w, jsonError("You are not the owner of this collection"), http.StatusForbidden)
		return
	}

	err = cc.CollectionUseCase.DeleteCard(r.Context(), collectionID, cardID)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Card deleted successfully",
	})
}
