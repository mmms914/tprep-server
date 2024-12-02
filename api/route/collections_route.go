package route

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/gookit/slog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"main/database"
	"main/models"
	"net/http"
	"strconv"
)

func getCollection(w http.ResponseWriter, r *http.Request) {
	var result models.Collection
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	filter := bson.D{{Key: "id", Value: id}}
	err := collections.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("collection not found"))
		} else {
			slog.FatalErr(err)
		}
		return
	}
	s, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(s)
}

func createCollection(w http.ResponseWriter, r *http.Request) {
	var newCollection models.Collection
	var collectionInfo models.CollectionInfo
	var collectionName string

	collectionName = r.URL.Query().Get("name")
	if collectionName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid name"))
		return
	}

	collectionTypeS := r.URL.Query().Get("is_public")
	collectionType, err := strconv.ParseBool(collectionTypeS)
	if collectionTypeS != "" && err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid is_public"))
		return
	}

	newCollection.Name = collectionName
	if collectionTypeS != "" {
		newCollection.IsPublic = collectionType
	}
	newCollection.Cards = make([]models.Card, 0)

	globalValues = database.GlobalValues()
	globalValues.MaxCollectionID += 1
	newCollection.ID = globalValues.MaxCollectionID
	database.SetGlobalValues(globalValues)

	_, err = collections.InsertOne(context.TODO(), newCollection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	collectionInfo.Id = newCollection.ID
	collectionInfo.MaxId = 1

	_, err = collectionsGlobal.InsertOne(context.TODO(), collectionInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response, _ := json.Marshal(newCollection)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func editCollection(w http.ResponseWriter, r *http.Request) {
	var newCollection, existingCollection models.Collection
	var collectionName string

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid id"))
		return
	}

	collectionName = r.URL.Query().Get("name")
	if r.URL.Query().Has("name") && collectionName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid name"))
		return
	}

	collectionTypeS := r.URL.Query().Get("is_public")
	collectionType, err := strconv.ParseBool(collectionTypeS)
	if collectionTypeS != "" && err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid is_public"))
		return
	}

	if !r.URL.Query().Has("name") && !r.URL.Query().Has("is_public") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid parameters"))
		return
	}

	filter := bson.D{{"id", id}}
	err = collections.FindOne(context.TODO(), filter).Decode(&existingCollection)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("collection not found"))
		return
	}

	if !r.URL.Query().Has("name") {
		collectionName = existingCollection.Name
	}

	if !r.URL.Query().Has("is_public") {
		collectionType = existingCollection.IsPublic
	}

	update := bson.D{
		{"$set", bson.D{
			{"name", collectionName},
			{"is_public", collectionType},
		}},
	}

	err = collections.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&newCollection)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("collection not found"))
		return
	}

	response, _ := json.Marshal(newCollection)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func deleteCollection(w http.ResponseWriter, r *http.Request) {
	var result models.Collection
	var resultInfo models.CollectionInfo
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	filter := bson.D{{Key: "id", Value: id}}
	err := collections.FindOneAndDelete(context.TODO(), filter).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("collection not found"))
		return
	}

	err = collectionsGlobal.FindOneAndDelete(context.TODO(), filter).Decode(&resultInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("collection successfully deleted"))
}

func initCollectionRouter(r *chi.Mux) {
	r.Get("/collection/{id}", getCollection)
	r.Post("/collection", createCollection)
	r.Put("/collection/{id}", editCollection)
	r.Delete("/collection/{id}", deleteCollection)
}
