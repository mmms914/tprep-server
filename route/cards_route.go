package route

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"main/models"
	"net/http"
	"strconv"
)

func addCard(w http.ResponseWriter, r *http.Request) {
	var card models.Card
	var collection models.Collection
	var collectionInfo models.CollectionInfo

	collectionId, err := strconv.Atoi(r.URL.Query().Get("collection_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid collection_id"))
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&card); err != nil || card.Question == "" || card.Answer == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid body data"))
		return
	}

	filter := bson.D{{Key: "id", Value: collectionId}}
	if err = collections.FindOne(context.TODO(), filter).Decode(&collection); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid collection"))
		return
	}

	err = collectionsGlobal.FindOne(context.TODO(), filter).Decode(&collectionInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	card.LocalID = collectionInfo.MaxId + 1

	collection.AddCard(card)
	_, err = collections.ReplaceOne(context.TODO(), filter, collection)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	collectionInfo.MaxId += 1
	_, err = collectionsGlobal.ReplaceOne(context.TODO(), filter, collectionInfo)

	s, _ := json.Marshal(collection)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(s))
}

func updateCard(w http.ResponseWriter, r *http.Request) {
	var collection models.Collection
	var card models.Card

	collectionId, err := strconv.Atoi(r.URL.Query().Get("collection_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid collection_id"))
		return
	}

	cardLocalId, err := strconv.Atoi(r.URL.Query().Get("card_local_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid card_local_id"))
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&card); err != nil || (card.Question == "" && card.Answer == "") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid body data"))
		return
	}

	filter := bson.D{{Key: "id", Value: collectionId}}
	if err = collections.FindOne(context.TODO(), filter).Decode(&collection); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid collection"))
		return
	}

	index, err := collection.UpdateCard(cardLocalId, card)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = collections.ReplaceOne(context.TODO(), filter, collection)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	s, _ := json.Marshal(collection.Cards[index])
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

func deleteCard(w http.ResponseWriter, r *http.Request) {
	var collection models.Collection

	collectionId, err := strconv.Atoi(r.URL.Query().Get("collection_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid collection_id"))
		return
	}

	cardLocalId, err := strconv.Atoi(r.URL.Query().Get("card_local_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid card_local_id"))
		return
	}

	filter := bson.D{{Key: "id", Value: collectionId}}
	if err = collections.FindOne(context.TODO(), filter).Decode(&collection); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid collection"))
		return
	}

	if collection.DeleteCard(cardLocalId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Card not found"))
		return
	}

	_, err = collections.ReplaceOne(context.TODO(), filter, collection)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Card was deleted"))
}

func initCardRouter(r *chi.Mux) {
	r.Post("/card", addCard)
	r.Put("/card", updateCard)
	r.Delete("/card", deleteCard)
}
