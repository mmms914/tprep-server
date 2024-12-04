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
	err = cc.CollectionUseCase.Create(r.Context(), &collection)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Collection created successfully",
	})
}

func (cc *CollectionController) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	collection, err := cc.CollectionUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(collection)
}

func (cc *CollectionController) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := cc.CollectionUseCase.DeleteByID(r.Context(), id)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Collection deleted successfully",
	})
}

func (cc *CollectionController) CreateCard(w http.ResponseWriter, r *http.Request) {
	var card domain.Card
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
	err = cc.CollectionUseCase.AddCard(r.Context(), collectionID, &card)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.SuccessResponse{
		Message: "Card created successfully",
	})
}

func (cc *CollectionController) DeleteCard(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	cardID, err := strconv.Atoi(chi.URLParam(r, "cardID"))
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
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
